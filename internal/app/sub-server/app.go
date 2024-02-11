package grpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/humanbelnik/backy-auth-ms/internal/grpc/auth"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

// In this package we manage gRPC server itself.
// We run with "Must" so any error while server launch will cause panic.
// Server shutting down graefully in order to complete all currenbt tasks.

type ServerApplication struct {
	Log    *slog.Logger
	Server *grpc.Server
	Port   int
}

func NewApplication(log *slog.Logger, auth authgrpc.Auth, port int) *ServerApplication {
	server := grpc.NewServer()
	authgrpc.Register(server, auth)

	return &ServerApplication{
		Log:    log,
		Server: server,
		Port:   port,
	}
}

func (s *ServerApplication) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *ServerApplication) Run() error {
	const fn = "grpcapp.Run"
	log := s.Log.With(
		slog.String("fn", fn),
		slog.Int("port", s.Port),
	)

	// TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("listening")

	// Run server
	if err := s.Server.Serve(listener); err != nil {
		return fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("serving")

	return nil
}

func (s *ServerApplication) Stop() {
	const fn = "grpcapp.Stop"
	log := s.Log.With(
		slog.String("fn", fn),
		slog.Int("port", s.Port),
	)
	log.Info("stopping server")

	// Stop
	s.Server.GracefulStop()
}
