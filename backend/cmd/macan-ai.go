package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/qwaq-dev/macan-ai/internal/config"
	"github.com/qwaq-dev/macan-ai/internal/handlers"
	"github.com/qwaq-dev/macan-ai/internal/repository/postgres"
	"github.com/qwaq-dev/macan-ai/internal/routes"
	"github.com/qwaq-dev/macan-ai/internal/services"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	app := fiber.New()
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	log.Info("Starting macan-ai", slog.String("env", cfg.Env))
	db, err := postgres.InitDatabase(cfg.Database, log)
	if err != nil {
		log.Error("Error with connecting to database", sl.Err(err))
		os.Exit(1)
	}

	userRepo := postgres.NewUserRepo(log, db)
	userService := services.NewUserService(log, userRepo)
	userHandler := handlers.NewUserHandler(log, userService)

	resumeParsingServiceAddr := cfg.Services.ResumeParsingGRPCAddr
	if resumeParsingServiceAddr == "" {
		log.Error("Resume Parsing gRPC address is empty in config. Check RESUME_PARSING_GRPC_ADDR environment variable or config file.",
			slog.String("config_path", os.Getenv("CONFIG_PATH")))
		os.Exit(1)
	}

	resumeParsingClient, err := services.NewResumeParsingClient(log, resumeParsingServiceAddr)
	if err != nil {
		log.Error("Failed to create Resume Parsing gRPC client", sl.Err(err))
		os.Exit(1)
	}
	defer func() {
		if err := resumeParsingClient.Close(); err != nil {
			log.Error("Error closing Resume Parsing gRPC client connection", sl.Err(err))
		} else {
			log.Info("Resume Parsing gRPC client connection closed.")
		}
	}()

	userResumeRepo := postgres.NewUserResumeRepo(log, db)
	userResumeService := services.NewUserResumeService(log, resumeParsingClient, userResumeRepo)
	userResumeHandler := handlers.NewUserResumeHandler(log, userResumeService)

	routes.InitRoutes(app, log, userHandler, userResumeHandler)

	log.Info("starting server", slog.String("address", cfg.Server.Port))

	go func() {
		if err := app.Listen(cfg.Server.Port); err != nil {
			log.Error("Fiber server failed to start", sl.Err(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Error with shutting down Fiber server", sl.Err(err))
	} else {
		log.Info("Fiber server gracefully stopped.")
	}

	log.Info("Application exited.")

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
