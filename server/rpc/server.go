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

package rpc

import (
	"context"
	"fmt"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/metis-labs/metis-server/api"
	"github.com/metis-labs/metis-server/api/converter"
	"github.com/metis-labs/metis-server/internal/log"
	"github.com/metis-labs/metis-server/server/database"
	"github.com/metis-labs/metis-server/server/projects"
	"github.com/metis-labs/metis-server/server/types"
	"github.com/metis-labs/metis-server/server/yorkie"
)

// Config is the configuration for creating a Server instance.
type Config struct {
	Port     int
	CertFile string
	KeyFile  string
}

// Server is a normal server that processes the logic requested by the client.
type Server struct {
	pb.UnimplementedMetisServer

	conf       *Config
	yorkieConf *yorkie.Config
	db         database.Database
	grpcServer *grpc.Server
}

// NewServer creates a new instance of Server.
func NewServer(conf *Config, yorkieConf *yorkie.Config, db database.Database) (*Server, error) {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			unaryInterceptor,
		)),
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			streamInterceptor,
		)),
	}

	if conf.CertFile != "" && conf.KeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(conf.CertFile, conf.KeyFile)
		if err != nil {
			log.Logger.Error(err)
			return nil, err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	rpcServer := &Server{
		conf:       conf,
		yorkieConf: yorkieConf,
		db:         db,
		grpcServer: grpc.NewServer(opts...),
	}
	pb.RegisterMetisServer(rpcServer.grpcServer, rpcServer)

	return rpcServer, nil
}

// Start starts to handle requests on incoming connections.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.conf.Port))
	if err != nil {
		return err
	}

	log.Logger.Infof("RPCServer is running on %d", s.conf.Port)

	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			log.Logger.Warnf("grpc server: %s", err.Error())
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
	project, err := projects.Create(ctx, s.db, s.yorkieConf, req.ProjectName)
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
	projectList, err := s.db.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListProjectsResponse{
		Projects: converter.ToProjects(projectList),
	}, nil
}

// UpdateProject updates the given project.
func (s *Server) UpdateProject(
	ctx context.Context,
	req *pb.UpdateProjectRequest,
) (*pb.UpdateProjectResponse, error) {
	if err := s.db.UpdateProject(ctx, types.ID(req.ProjectId), req.ProjectName); err != nil {
		return nil, err
	}

	return &pb.UpdateProjectResponse{}, nil
}

// DeleteProject deletes the given project.
func (s *Server) DeleteProject(
	ctx context.Context,
	req *pb.DeleteProjectRequest,
) (*pb.DeleteProjectResponse, error) {
	if err := s.db.DeleteProject(ctx, types.ID(req.ProjectId)); err != nil {
		return nil, err
	}

	return &pb.DeleteProjectResponse{}, nil
}
