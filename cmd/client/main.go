package main

import (
	"client/internal/app"
	"client/internal/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

	go client_application.Serv.RunServer()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	logger.Info("stopping the application")
	client_application.Serv.Stop()
	logger.Info("application successfully stoped")
}

func createLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
