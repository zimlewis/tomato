package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pkgerrs "github.com/pkg/errors"
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

	port := os.Getenv("TOMATO_PORT")
	if port == "" { port = "6600" }
	addr := fmt.Sprintf("0.0.0.0:%s", port)


	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Cannot start server: %w", pkgerrs.WithStack(err))
	}
	defer func(){
		if err := listener.Close(); err!= nil {
			fmt.Println("error closing server: ", err)
			return
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger)),
		),
	)

	repo := badgerrepo.New(storage.Storage)
	service := timer.New(&repo, logger)

	proto.RegisterTimerServer(grpcServer, service)

	errChan := make(chan error, 1)

	go func() {
		errChan <- grpcServer.Serve(listener)
	}()
	logger.Debug("serving", slog.String("port", "6601"))

	select {
		case err := <- errChan: return err
		case <- ctx.Done(): return nil
	}
}
