package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Connection *grpc.ClientConn
}

func New() (*Client, error) {
	var client Client;
	conn, err := grpc.NewClient(
		"dns:///localhost:6600", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client.Connection = conn

	return &client, nil
}
