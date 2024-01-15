package main

import (
	"client/internal/app"
	"client/internal/config"
	"context"
	"log/slog"
	"os"
)

func main() {
	logger := createLogger()
	if logger == nil {
		panic("logger failed to create")
	}

	logger.Info("logger has been created successfully")

	cfg, err := config.ReadConfig()
	if err != nil {
		logger.Error("reading config error", slog.Any("err", err))
		return
	}

	logger.Info("config has been read successfully")

	client_application, err := app.NewApp(logger, cfg.Http, context.Background(), cfg.Grpc)
	if err != nil {
		logger.Error("creation client_application error", slog.Any("err", err))
		return
	}
	client_application.Serv.RunServer()

	// TODO: graceful stop
}

func createLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
