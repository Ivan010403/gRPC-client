package app

import (
	grpcclient "client/internal/app/grpcClient"
	httpserver "client/internal/app/http-server"
	"client/internal/config"
	"context"
	"log/slog"
)

type App struct {
	Serv httpserver.HTTP_server
}

func NewApp(log *slog.Logger, cfg config.HTTPServer, ctx context.Context, client config.GRPC_server) (*App, error) {
	//TODO: защищенное соединение, шифрование
	//TODO: закрывать соединение
	clnt, err := grpcclient.NewClient(log, ctx, client.Address, client.Timeout, client.Max_retry)
	if err != nil {
		return nil, err
	}
	log.Info("creation of gRPC-client is successfull")

	serv := httpserver.NewServer(log, cfg.Address, cfg.Timeout, cfg.IdleTimeout, clnt)
	log.Info("creation of HTTP-server is successfull")

	return &App{Serv: *serv}, nil
}
