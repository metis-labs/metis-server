package client

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "oss.navercorp.com/metis/metis-server/api"
)

const (
	rpcAddr = "localhost:10118"
	timeout = 10 * time.Second
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MetisClient
}

func New() (*Client, error) {
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewMetisClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

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

func (c *Client) ListProjects(ctx context.Context) ([]*pb.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := c.client.ListProjects(ctx, &pb.ListProjectsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Projects, nil
}

func (c *Client) UpdateProject(ctx context.Context, projectID string, projectName string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := c.client.UpdateProject(ctx, &pb.UpdateProjectRequest{
		ProjectId:   projectID,
		ProjectName: projectName,
	})
	return err
}

func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := c.client.DeleteProject(ctx, &pb.DeleteProjectRequest{
		ProjectId: projectID,
	})
	return err
}
