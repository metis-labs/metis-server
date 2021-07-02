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

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/metis-labs/metis-server/api"
)

const (
	timeout = 10 * time.Second
)

// Client is a normal client that can communicate with the Server.
type Client struct {
	conn   *grpc.ClientConn
	client pb.MetisClient
}

// Option configures how we set up the client.
type Option struct {
	UserID   string
	CertFile string
}

// Dial creates an instance of Client.
func Dial(rpcAddr string, opts ...Option) (*Client, error) {
	var userID string
	if len(opts) > 0 && opts[0].UserID != "" {
		userID = opts[0].UserID
	}

	conn, err := grpc.Dial(
		rpcAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryInterceptor(userID)),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewMetisClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes all resources of this client.
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateProject creates a new client of the given name.
func (c *Client) CreateProject(ctx context.Context, name string) (*pb.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := c.client.CreateProject(ctx, &pb.CreateProjectRequest{
		ProjectName: name,
	})
	if err != nil {
		return nil, err
	}

	return res.Project, nil
}

// ListProjects returns the list of clients.
func (c *Client) ListProjects(ctx context.Context) ([]*pb.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := c.client.ListProjects(ctx, &pb.ListProjectsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Projects, nil
}

// UpdateProject updates the given project.
func (c *Client) UpdateProject(ctx context.Context, projectID string, projectName string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := c.client.UpdateProject(ctx, &pb.UpdateProjectRequest{
		ProjectId:   projectID,
		ProjectName: projectName,
	})
	return err
}

// DeleteProject deletes the given project.
func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := c.client.DeleteProject(ctx, &pb.DeleteProjectRequest{
		ProjectId: projectID,
	})
	return err
}
