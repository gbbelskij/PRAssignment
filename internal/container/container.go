package container

import (
	"PRAssignment/internal/config"
	"PRAssignment/internal/logger"
	"PRAssignment/internal/repository/storage"
	"context"
	"log/slog"
	"os"
)

type Container struct {
	Config  *config.Config
	Logger  *slog.Logger
	Storage *storage.Storage
}

func NewContainer() *Container {
	cfg := config.MustLoad()
	log := logger.NewLogger(cfg.Env)
	pgstorage, err := storage.NewStorage(context.Background())
	if err != nil {
		log.Error("failed to connect to database")
		os.Exit(2)
	}
	return &Container{
		Config:  cfg,
		Logger:  log,
		Storage: pgstorage,
	}
}
