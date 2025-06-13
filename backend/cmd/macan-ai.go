package main

import (
	"log/slog"
	"os"

	"github.com/qwaq-dev/macan-ai/internal/config"
	"github.com/qwaq-dev/macan-ai/internal/repository/postgres"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("Starting macan-ai", slog.String("env", cfg.Env))
	database, err := postgres.InitDatabase(cfg.Database, log)
	if err != nil {
		log.Error("Error with connecting to database", sl.Err(err))
	}

	_ = database

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
