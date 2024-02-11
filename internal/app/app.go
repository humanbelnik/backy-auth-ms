package app

import (
	"time"

	grpcapp "github.com/humanbelnik/backy-auth-ms/internal/app/sub-server"
	"github.com/humanbelnik/backy-auth-ms/internal/config"
	auth_service "github.com/humanbelnik/backy-auth-ms/internal/services/auth"
	"github.com/humanbelnik/backy-auth-ms/internal/storage/postgres"
	"golang.org/x/exp/slog"
)

type MainApplication struct {
	Server *grpcapp.ServerApplication
}

// NewMainApplication creates new application with gRPC-server, Storage and Auth service.
func NewMainApplication(log *slog.Logger, port int, databaseConfig config.DatabaseConfig, tokenTTL time.Duration) *MainApplication {
	const fn = "app.NewMainApplication"
	log = log.With(
		slog.String("fn", fn),
	)

	storage, err := postgres.New(log, databaseConfig)
	if err != nil {
		panic(err)
	}
	log.Info("storage created")

	authService := auth_service.New(log, storage, storage, tokenTTL)
	grpcApp := grpcapp.NewApplication(log, authService, port)

	return &MainApplication{
		Server: grpcApp,
	}
}
