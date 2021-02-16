package server

import (
	"oss.navercorp.com/metis/metis-server/server/api"
)

type Server struct {
	rpcServer *api.RPCServer

	shutdown   bool
	shutdownCh chan struct{}
}

const rpcPort = 10118

func New() (*Server, error) {
	rpcServer, err := api.NewRPCServer()
	if err != nil {
		return nil, err
	}

	return &Server{
		rpcServer:  rpcServer,
		shutdownCh: make(chan struct{}),
	}, nil
}

func (s *Server) Start() error {
	return s.rpcServer.Start(rpcPort)
}

func (s *Server) Shutdown(graceful bool) error {
	if s.shutdown {
		return nil
	}

	if graceful {
		s.rpcServer.GracefulStop()
	} else {
		s.rpcServer.Stop()
	}

	s.shutdown = true
	close(s.shutdownCh)

	return nil
}

func (s *Server) ShutdownCh() <-chan struct{} {
	return s.shutdownCh
}
