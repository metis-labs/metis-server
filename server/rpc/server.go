package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "oss.navercorp.com/metis/metis-server/api"
	"oss.navercorp.com/metis/metis-server/api/converter"
	"oss.navercorp.com/metis/metis-server/internal/log"
	"oss.navercorp.com/metis/metis-server/server/database"
)

// Server is a normal server that processes the logic requested by the client.
type Server struct {
	db         database.Database
	grpcServer *grpc.Server
}

// NewServer creates a new instance of Server.
func NewServer(db database.Database) (*Server, error) {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			unaryInterceptor,
		)),
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			streamInterceptor,
		)),
	}

	rpcServer := &Server{
		db:         db,
		grpcServer: grpc.NewServer(opts...),
	}
	pb.RegisterMetisServer(rpcServer.grpcServer, rpcServer)

	return rpcServer, nil
}

// Start starts to handle requests on incoming connections.
func (s *Server) Start(rpcPort int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", rpcPort))
	if err != nil {
		return err
	}

	log.Logger.Infof("RPCServer is running on %d", rpcPort)

	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			log.Logger.Errorf("fail to serve: %s", err.Error())
		} else {
			log.Logger.Info("grpc server closed")
		}
	}()
	return nil
}

// GracefulStop stops the gRPC server gracefully.
func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}

// Stop stops the gRPC server. It immediately closes all open
// connections and listeners.
func (s *Server) Stop() {
	s.grpcServer.Stop()
}

// CreateProject creates a new project of the given name.
func (s *Server) CreateProject(
	ctx context.Context,
	req *pb.CreateProjectRequest,
) (*pb.CreateProjectResponse, error) {
	project, err := s.db.CreateProject(ctx, req.ProjectName)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProjectResponse{
		Project: converter.ToProject(project),
	}, nil
}

// ListProjects returns the list of projects.
func (s *Server) ListProjects(
	ctx context.Context,
	req *pb.ListProjectsRequest,
) (*pb.ListProjectsResponse, error) {
	projects, err := s.db.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListProjectsResponse{
		Projects: converter.ToProjects(projects),
	}, nil
}

// UpdateProject updates the given project.
func (s *Server) UpdateProject(
	ctx context.Context,
	req *pb.UpdateProjectRequest,
) (*pb.UpdateProjectResponse, error) {
	if err := s.db.UpdateProject(ctx, req.ProjectId, req.ProjectName); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else if errors.Is(err, database.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, err
	}

	return &pb.UpdateProjectResponse{}, nil
}

// DeleteProject deletes the given project.
func (s *Server) DeleteProject(
	ctx context.Context,
	req *pb.DeleteProjectRequest,
) (*pb.DeleteProjectResponse, error) {
	if err := s.db.DeleteProject(ctx, req.ProjectId); err != nil {
		return nil, err
	}

	return &pb.DeleteProjectResponse{}, nil
}
