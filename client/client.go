package client

import (
	"fmt"

	"google.golang.org/grpc"
)

// A gRPC client that hold the connection to the gRPC server
type Client struct {
	port        string
	host        string
	dialOptions []grpc.DialOption
}
type option func(*Client)

func (client Client) GetConnection() (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%s", client.host, client.port)
	conn, err := grpc.NewClient(
		addr,
		client.dialOptions...,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func WithHost(host string) option {
	return func(c *Client) {
		c.host = host
	}
}

func WithPort(port string) option {
	return func(c *Client) {
		c.port = port
	}
}

func WithDialOptions(options ...grpc.DialOption) option {
	return func(c *Client) {
		c.dialOptions = append(c.dialOptions, options...)
	}
}

// New gRPC client that connect to the configured port of tomato
func New(opts ...option) *Client {
	client := new(Client{
		port:        "6600",
		host:        "dns:///localhost",
		dialOptions: []grpc.DialOption{},
	})

	for _, opt := range opts {
		opt(client)
	}

	return client
}
