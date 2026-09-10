package client

import (
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Connection *grpc.ClientConn
}

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
