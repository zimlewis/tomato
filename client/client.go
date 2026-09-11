package client

import (
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// A gRPC client that hold the connection to the gRPC server
type Client struct {
	Connection *grpc.ClientConn
}

// New gRPC client that connect to the configured port of tomato
func New() (*Client, error) {
	port := os.Getenv("TOMATO_PORT")
	if port == "" {
		port = "6600"
	}

	addr := fmt.Sprintf("dns:///localhost:%s", port)
	conn, err := grpc.NewClient(
		addr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := new(Client)
	client.Connection = conn

	return client, nil
}
