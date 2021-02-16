package rpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "oss.navercorp.com/metis/metis-server/api"
)

type Server struct {
	grpcServer *grpc.Server
}

func (s *Server) CreateStudy(_ context.Context, req *pb.CreateStudyRequest) (*pb.CreateStudyResponse, error) {
	return &pb.CreateStudyResponse{
		Study: &pb.Study{
			Name: req.StudyName,
		},
	}, nil
}

func NewServer() (*Server, error) {
	rpcServer := &Server{
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
