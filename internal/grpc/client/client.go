package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/bxrne/goforge/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

type Config struct {
	Address string
	Timeout time.Duration
}

func NewClient(config Config) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetUser(ctx context.Context, id int32) (*pb.User, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (c *Client) CreateUser(ctx context.Context, name, email string) (*pb.User, error) {
	resp, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (c *Client) ListUsers(ctx context.Context, page, limit int32) ([]*pb.User, int32, error) {
	resp, err := c.client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.Users, resp.Total, nil
}
