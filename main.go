package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/cmorent/killteam-wcw-rankings/pkg/db/storage"
	"github.com/cmorent/killteam-wcw-rankings/pkg/server"
)

type Config struct {
	Port         string `envconfig:"PORT" default:"8080"`
	DBBucketName string `envconfig:"DB_BUCKET_NAME" default:"kf-kt-wcw-rankings"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	ctx := context.Background()

	db, err := storage.New(ctx, cfg.DBBucketName)
	if err != nil {
		slog.Error("failed to init db", slog.Any("error", err))
		os.Exit(1)
	}

	srv, err := server.New(net.JoinHostPort("", cfg.Port), db)
	if err != nil {
		slog.Error("failed to init server", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := srv.Run(ctx); err != nil {
		slog.Error("failed to run the server", slog.Any("error", err))
		os.Exit(1)
	}
}
