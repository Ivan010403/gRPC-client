package app

import (
	"context"

	proto "github.com/Ivan010403/proto/protoc/go"
	"google.golang.org/grpc"
)

type Client struct {
	Client proto.CloudClient
}

func New(ctx context.Context, addr string) (*Client, error) {
	//TODO: защищенное соединение, шифрование
	//TODO: закрывать соединение
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{Client: proto.NewCloudClient(conn)}, nil
}
