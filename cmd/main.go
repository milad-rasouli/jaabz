package main

import (
	"context"
	"github.com/milad-rasouli/jaabz/internal/infra/godotenv"
	"github.com/milad-rasouli/jaabz/internal/infra/redis"
	"github.com/milad-rasouli/jaabz/internal/repo/duplicate"
	jaabz2 "github.com/milad-rasouli/jaabz/internal/repo/jaabz"
	"github.com/milad-rasouli/jaabz/internal/service"
	"log/slog"
	"os"
	"time"
)

// Job represents a job listing with relevant details

func main() {
	env := godotenv.NewEnv()
	logger := initSlogLogger()
	logger.Info("welcome to " + env.AppName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rdis := redis.NewRedis(env)
	err := rdis.Setup(ctx)
	if err != nil {
		logger.Error("failed to setup redis", "error", err)
		os.Exit(1)
	}
	defer rdis.Close()
	err = rdis.HealthCheck(ctx)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	dupl := duplicate.New(logger, rdis)
	jaabz := jaabz2.New(env, logger)

	jaabzService := service.NewJaabzService(logger, dupl, jaabz)
	err = jaabzService.JaabzProcess()
	if err != nil {
		logger.Error("failed to process Jaabz", "error", err)
		os.Exit(1)
	}

}

func initSlogLogger() *slog.Logger {
	slogHandlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, slogHandlerOptions))
	slog.SetDefault(logger)

	return logger
}
