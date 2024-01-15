package delete

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
		log.Info("starting delete http handler")

		name, format, err := validate(r.FormValue("namedelete"))
		if err != nil {
			log.Error("validate error", slog.Any("err", err))
			http.Error(w, "can't get name of file", http.StatusInternalServerError)
			return
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

		data, err := client.DeleteFile(ctx, name, format)
		if err != nil {
			log.Error("client DeleteFileErr", slog.Any("err", err))
			http.Error(w, "can't delete file from cloud", http.StatusInternalServerError)
			return
		}

		log.Info("file was deleted", slog.String("name", data))

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
