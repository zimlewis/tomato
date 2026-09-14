package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pkgerrs "github.com/pkg/errors"
	prototimer "github.com/zimlewis/tomato/gen/proto/timer"
	"google.golang.org/grpc"
)

type Server struct {
	port    string
	host    string
	logger  *slog.Logger
	service prototimer.TimerServer
}

type option func (*Server)

func WithPort(port string) option {
	return func(s *Server) {
		s.port = port
	}
}

func WithHost(host string) option {
	return func(s *Server) {
		s.host = host
	}
}

func WithLogger(logger *slog.Logger) option {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithService(service prototimer.TimerServer) option {
	return func(s *Server) {
		s.service = service
	}
}

func New(opts ...option) Server {
	server := Server{}

	for _, opt := range opts {
		opt(&server)
	}

	return server
}

func (server Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", server.host, server.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("cannot start server: %w", pkgerrs.WithStack(err))
	}
	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Println("error closing server: ", err)
			return
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(server.logger)),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptorLogger(server.logger)),
		),
	)

	prototimer.RegisterTimerServer(grpcServer, server.service)
	
	// Create an error channel to listen to grpc in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- grpcServer.Serve(listener)
	}()
	server.logger.Debug("serving", slog.String("port", server.port))

	select {
	case err := <-errChan: return err
	case <-ctx.Done(): return nil
	}
}

