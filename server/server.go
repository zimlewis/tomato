package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pkgerrs "github.com/pkg/errors"
	prototimer "github.com/zimlewis/tomato/gen/proto/timer"
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
	defer func() {
		if err := c(); err != nil {
			fmt.Printf("failed to flush logger: %v", err)
		}
	}()

	// Get the port, default to 6600
	port := os.Getenv("TOMATO_PORT")
	if port == "" {
		port = "6600"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	// Open a tcp listener to listen the grpc server on
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Cannot start server: %w", pkgerrs.WithStack(err))
	}
	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Println("error closing server: ", err)
			return
		}
	}()

	// Initialize a server that have a logger
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger)),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptorLogger(logger)),
		),
	)

	// Create new repo and asign it to the service along with the logger
	repo := badgerrepo.New(storage.Storage)
	service := timer.New(&repo, logger)

	// Register the timer service to the server
	prototimer.RegisterTimerServer(grpcServer, service)

	// Create an error channel to listen to grpc in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- grpcServer.Serve(listener)
	}()
	logger.Debug("serving", slog.String("port", port))

	select {
	case err := <-errChan: return err
	case <-ctx.Done(): return nil
	}
}
