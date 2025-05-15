package main

import (
	"context"
	"os"
	"time"

	"github.com/milad-rasouli/jaabz/internal/infra/godotenv"
	"github.com/milad-rasouli/jaabz/internal/infra/redis"
	"github.com/milad-rasouli/jaabz/internal/repo/duplicate"
	jaabz2 "github.com/milad-rasouli/jaabz/internal/repo/jaabz"
	"github.com/milad-rasouli/jaabz/internal/repo/telegram"
	"github.com/milad-rasouli/jaabz/internal/service"
	"log/slog"
)

func main() {
	logger := initSlogLogger()
	env := godotenv.NewEnv()
	logger.Info("Welcome to " + env.AppName)

	// Verify environment variables
	if env.TelegramBotToken == "" || env.TelegramChannelID == "" {
		logger.Error("Missing Telegram environment variables", "token_empty", env.TelegramBotToken == "", "channel_id_empty", env.TelegramChannelID == "")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rdis := redis.NewRedis(env)
	if err := rdis.Setup(ctx); err != nil {
		logger.Error("Failed to setup redis", "error", err)
		os.Exit(1)
	}
	defer rdis.Close()

	if err := rdis.HealthCheck(ctx); err != nil {
		logger.Error("Failed to connect to redis", "error", err)
		os.Exit(1)
	}

	dupl := duplicate.New(logger, rdis)
	jaabz := jaabz2.New(env, logger)
	tele, err := telegram.New(logger, env)
	if err != nil {
		logger.Error("Failed to initialize Telegram", "error", err)
		os.Exit(1)
	}

	jaabzService := service.NewJaabzService(logger, dupl, jaabz, tele)
	if err := jaabzService.StartJaabzProcess(context.Background()); err != nil {
		logger.Error("Failed to process Jaabz", "error", err)
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
