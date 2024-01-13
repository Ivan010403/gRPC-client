package main

import (
	"client/internal/app"
	"context"
	"log"

	proto "github.com/Ivan010403/proto/protoc/go"
)

func main() {

	parentContext := context.Background()
	ctx, _ := context.WithCancel(parentContext)

	app, err := app.New(ctx, ":4545")
	if err != nil {
		panic("app")
	}

	err = printRecord(app.Client, &proto.UploadFileRequest{NameFile: "gRPC Stream Client: Record", File: []byte{0, 23, 45, 67, 79, 65, 67}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}
}

func printRecord(client proto.CloudClient, r *proto.UploadFileRequest) error {
	stream, err := client.UploadFile(context.Background())
	if err != nil {
		return err
	}

	err = stream.Send(r)
	if err != nil {
		return err
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("resp: pj.name: %s", resp.GetNameFile())

	return nil
}
