package server

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	proto "github.com/zimlewis/tomato/gen/proto"
	"github.com/zimlewis/tomato/internal/badgerrepo"
	"github.com/zimlewis/tomato/internal/service/timer"
	"github.com/zimlewis/tomato/storage"
	"google.golang.org/grpc"
)


func Start(ctx context.Context) error {
	logger, c, err := initializeLogger()
	if err != nil {
		return fmt.Errorf("Cannot initialize logger: %v", err)
	}
	defer func () {
		if err := c(); err != nil {
			fmt.Printf("failed to flush logger: %v", err)
		}
	}()


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

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger)),
		),
	)

	repo := badgerrepo.New(storage.Storage)
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
