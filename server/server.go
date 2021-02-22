package server

import (
	"context"
	"log"

	"oss.navercorp.com/metis/metis-server/server/database"
	"oss.navercorp.com/metis/metis-server/server/database/mongodb"
	"oss.navercorp.com/metis/metis-server/server/rpc"
)

// Server receives requests from the client, stores data in the database,
type Server struct {
	rpcServer *rpc.Server
	db        database.Database

	shutdown   bool
	shutdownCh chan struct{}
}

const rpcPort = 10118

// New creates a new instance of Server.
func New() (*Server, error) {
	dbClient := mongodb.NewClient()
	rpcServer, err := rpc.NewServer(dbClient)
	if err != nil {
		return nil, err
	}

	return &Server{
		rpcServer:  rpcServer,
		db:         dbClient,
		shutdownCh: make(chan struct{}),
	}, nil
}

// Start starts the server by opening the rpc port.
func (s *Server) Start() error {
	if err := s.db.Dial(context.Background()); err != nil {
		return err
	}

	return s.rpcServer.Start(rpcPort)
}

// Shutdown shuts down this server.
func (s *Server) Shutdown(graceful bool) error {
	if s.shutdown {
		return nil
	}

	if graceful {
		s.rpcServer.GracefulStop()
	} else {
		s.rpcServer.Stop()
	}

	if err := s.db.Close(context.Background()); err != nil {
		log.Print(err)
	}

	s.shutdown = true
	close(s.shutdownCh)

	return nil
}

// ShutdownCh returns the shutdown channel.
func (s *Server) ShutdownCh() <-chan struct{} {
	return s.shutdownCh
}
