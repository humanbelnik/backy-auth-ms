package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/humanbelnik/backy-auth-ms/internal/app"
	"github.com/humanbelnik/backy-auth-ms/internal/config"
	"github.com/humanbelnik/logit/logit"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(".env not found")
	}
}

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	// Init application.
	// Run application in a sepereate goroutine in order to implement Gracefull shutdown.
	mainApplication := app.NewMainApplication(log, cfg.GRPC.Port, cfg.Database, cfg.TokenTTL)
	go mainApplication.Server.MustRun()

	// Shut down.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	mainApplication.Server.Stop()
	log.Info("main application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(logit.NewHandler(slog.LevelDebug))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
