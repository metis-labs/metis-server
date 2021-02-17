package rpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	pb "oss.navercorp.com/metis/metis-server/api"
	"oss.navercorp.com/metis/metis-server/server/database"
)

type Server struct {
	db         database.Database
	grpcServer *grpc.Server
}

func (s *Server) CreateModel(
	ctx context.Context,
	req *pb.CreateModelRequest,
) (*pb.CreateModelResponse, error) {
	model, err := s.db.CreateModel(ctx, req.ModelName)
	if err != nil {
		return nil, err
	}

	return &pb.CreateModelResponse{
		Model: &pb.Model{
			Name: model.Name,
		},
	}, nil
}

func (s *Server) CreateProject(
	ctx context.Context,
	req *pb.CreateProjectRequest,
) (*pb.CreateProjectResponse, error) {
	_, err := s.db.CreateProject(ctx, req.ProjectName)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProjectResponse{
		// Project: converter.toProject(project),
	}, nil

}

func (s *Server) ListProjects(
	ctx context.Context,
	req *pb.ListProjectsRequest,
) (*pb.ListProjectsResponse, error) {
	_, err := s.db.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListProjectsResponse{
		// Projects: converter.toProjects(projects),
	}, nil
}

func NewServer(db database.Database) (*Server, error) {
	rpcServer := &Server{
		db:         db,
		grpcServer: grpc.NewServer(),
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

	fmt.Printf("RPCServer is running on %d", rpcPort)

	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			fmt.Printf("fail to serve: %s", err.Error())
		} else {
			fmt.Printf("grpc server closed")
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
