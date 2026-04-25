package server

import (
	"context"
	"fmt"
	"net"

	proto "github.com/zimlewis/tomato/gen/proto"
	"github.com/zimlewis/tomato/internal/repository"
	"github.com/zimlewis/tomato/internal/service/timer"
	"github.com/zimlewis/tomato/storage"
	"google.golang.org/grpc"
)

func Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", "localhost:6600")
	defer func(){
		if err := listener.Close(); err!= nil {
			fmt.Println("error closing server: ", err)
			return
		}
	}()
	if err != nil {
		return fmt.Errorf("Cannot start server: %w", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	repo := repository.New(storage.Storage)
	service := timer.New(&repo)

	proto.RegisterTimerServer(grpcServer, service)

	errChan := make(chan error, 1)

	go func() {
		errChan <- grpcServer.Serve(listener)
	}()

	select {
		case err := <- errChan: return err
		case <- ctx.Done(): return nil
	}
}
