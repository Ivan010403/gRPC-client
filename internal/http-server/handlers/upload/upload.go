package upload

import (
	grpcclient "client/internal/app/grpcClient"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const size = 1024 * 1024 * 10

func New(log *slog.Logger, client *grpcclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("starting upload http handler")

		r.Body = http.MaxBytesReader(w, r.Body, size)
		if err := r.ParseMultipartForm(size); err != nil {
			http.Error(w, "The uploaded file is too big. Please choose an file that's less than 10MB in size", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name, format, err := getNameWithFormat(header.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

		full_name, err := client.UploadFile(context.Background(), data, name, format)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("file was uploaded", slog.String("name", full_name))
		//TODO: write response
	}
}

func getNameWithFormat(full_name string) (string, string, error) {
	for i := len(full_name) - 1; i >= 0; i-- {
		if rune(full_name[i]) == rune('.') {
			return full_name[0:i], full_name[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("finding format file error")
}
