package get

import (
	grpcclient "client/internal/app/grpcClient"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger, client *grpcclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("starting Get http handler")

		name, format, err := validate(r.FormValue("name"))
		if err != nil {
			log.Error("validation error", slog.Any("err", err))
			http.Error(w, "can't get name of file", http.StatusInternalServerError)
			return
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

		data, err := client.GetFile(ctx, name, format)
		if err != nil {
			log.Error("client getting file error", slog.Any("err", err))
			http.Error(w, "can't getfile from cloud", http.StatusInternalServerError)
			return
		}

		log.Info("file was gotten")

		_, err = w.Write(data)
		if err != nil {
			log.Error("writng response error", slog.Any("err", err))
			http.Error(w, "can't rescpond from cloud", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func validate(full_name string) (string, string, error) {
	var end_file int = 0

	for i := len(full_name) - 1; i >= 0; i-- {
		if (rune(full_name[i]) != ' ') && (rune(full_name[i]) != '\t') {
			if i > end_file {
				end_file = i
			}
		}
		if rune(full_name[i]) == rune('.') {
			return full_name[0:i], full_name[i+1 : end_file+1], nil
		}
	}
	return "", "", fmt.Errorf("finding format file error")
}
