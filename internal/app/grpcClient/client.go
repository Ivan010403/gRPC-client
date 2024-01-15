package grpcclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	proto "github.com/Ivan010403/proto/protoc/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	CldClient proto.CloudClient
	log       *slog.Logger
}

const chunkSize = 1024

func NewClient(logger *slog.Logger, ctx context.Context, addr string, timeout time.Duration, max_retry int) (*Client, error) {
	//TODO: защищенное соединение, шифрование
	rtrOpt := []retry.CallOption{
		retry.WithCodes(codes.NotFound, codes.DeadlineExceeded),
		retry.WithMax(uint(max_retry)),
		retry.WithPerRetryTimeout(timeout),
	}

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(retry.UnaryClientInterceptor(rtrOpt...)))
	if err != nil {
		return nil, fmt.Errorf("connection creation error: %w", err)
	}

	return &Client{CldClient: proto.NewCloudClient(conn), log: logger}, nil
}

func (c *Client) UploadFile(ctx context.Context, data []byte, name, format_file string) (string, error) {
	c.log.Info("starting UploadFile client")

	stream, err := c.CldClient.UploadFile(ctx)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(data)
	buffer := make([]byte, chunkSize)

	req := &proto.UploadFileRequest{NameFile: name, FileFormat: format_file}

	err = stream.Send(req)
	if err != nil {
		return "", err
	}

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		req := &proto.UploadFileRequest{File: buffer[:n]}

		err = stream.Send(req)

		if err != nil {
			return "", err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	c.log.Info("completed UploadFile client")

	return resp.GetFullName(), nil
}

func (c *Client) DeleteFile(ctx context.Context, name, format_file string) (string, error) {
	c.log.Info("starting DeleteFile client")

	resp, err := c.CldClient.DeleteFile(ctx, &proto.DeleteFileRequest{NameFile: name, FileFormat: format_file})

	if err != nil {
		return "", err
	}

	c.log.Info("completed DeleteFile client")
	return resp.GetFullName(), nil
}

func (c *Client) GetFile(ctx context.Context, name, format_file string) ([]byte, error) {
	c.log.Info("starting GetFile client")

	stream, err := c.CldClient.GetFile(ctx, &proto.GetFileRequest{NameFile: name, FileFormat: format_file})
	if err != nil {
		return nil, err
	}

	file_bytes := bytes.Buffer{}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		_, err = file_bytes.Write(resp.GetFile())
		if err != nil {
			return nil, err
		}
	}

	c.log.Info("completed GetFile client")
	return file_bytes.Bytes(), nil
}

func (c *Client) GetFullData(ctx context.Context) ([]struct {
	Name          string
	Creation_date string
	Update_date   string
}, error) {
	c.log.Info("starting GetFullData client")

	stream, err := c.CldClient.GetFullData(ctx, &proto.GetFullDataRequest{Id: "0"})
	if err != nil {
		return nil, err
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	data := make([]struct {
		Name          string
		Creation_date string
		Update_date   string
	}, resp.GetSize())

	counter := 0

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data[counter].Name = resp.GetName()
		data[counter].Creation_date = resp.GetCreationDate()
		data[counter].Update_date = resp.GetUpdatingDate()
		counter++
	}

	c.log.Info("completed GetFullData client")
	return data, nil
}
