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

func (s *Server) CreateDiagram(
	ctx context.Context,
	req *pb.CreateDiagramRequest,
) (*pb.CreateDiagramResponse, error) {
	diagram, err := s.db.CreateDiagram(ctx, req.DiagramName)
	if err != nil {
		return nil, err
	}

	return &pb.CreateDiagramResponse{
		Diagram: &pb.Diagram{
			Name: diagram.Name,
		},
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

// `Start` starts to handle requests on incoming connections.
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
