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

func (c *Client) CreateDiagram(ctx context.Context, diagramName string) (*pb.Diagram, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := c.client.CreateDiagram(ctx, &pb.CreateDiagramRequest{
		DiagramName: diagramName,
	})
	if err != nil {
		return nil, err
	}

	return res.Diagram, nil
}
