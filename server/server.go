/*
 * Copyright 2021-present NAVER Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"log"

	"github.com/metis-labs/metis-server/server/database"
	"github.com/metis-labs/metis-server/server/database/mongodb"
	"github.com/metis-labs/metis-server/server/rpc"
	"github.com/metis-labs/metis-server/server/web"
)

// Server receives requests from the client, stores data in the database,
type Server struct {
	conf      *Config
	rpcServer *rpc.Server
	webServer *web.Server
	db        database.Database

	shutdown   bool
	shutdownCh chan struct{}
}

// New creates a new instance of Server.
func New(conf *Config) (*Server, error) {
	dbClient := mongodb.NewClient(conf.Mongo)
	rpcServer, err := rpc.NewServer(conf.RPC, conf.Yorkie, dbClient)
	if err != nil {
		return nil, err
	}

	webServer, err := web.NewServer(conf.Web, dbClient, conf.Yorkie)
	if err != nil {
		return nil, err
	}

	return &Server{
		conf:       conf,
		rpcServer:  rpcServer,
		webServer:  webServer,
		db:         dbClient,
		shutdownCh: make(chan struct{}),
	}, nil
}

// Start starts the server by opening the rpc port.
func (s *Server) Start() error {
	if err := s.db.Dial(context.Background()); err != nil {
		return err
	}

	if err := s.webServer.Start(); err != nil {
		return err
	}

	return s.rpcServer.Start()
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

	if graceful {
		s.webServer.GracefulStop()
	} else {
		s.webServer.Stop()
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

// RPCAddr returns the RPC address.
func (s *Server) RPCAddr() string {
	return s.conf.RPCAddr()
}
