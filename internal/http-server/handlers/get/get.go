package get

import (
	grpcclient "client/internal/app/grpcClient"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func New(log *slog.Logger, client *grpcclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("starting Get http handler")

		name, format, err := validate(r.FormValue("name"))
		if err != nil {
			http.Error(w, "can't get name of file", http.StatusInternalServerError)
			return
		}

		data, err := client.GetFile(context.Background(), name, format)
		if err != nil {
			http.Error(w, "can't getfile from cloud", http.StatusInternalServerError)
			return
		}

		fl, err := os.Create("../../internal/web-site/public/static/img/" + "temp." + format)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "can't getfile from cloud", http.StatusInternalServerError)
			return
		}

		_, err = fl.Write(data)
		if err != nil {
			http.Error(w, "can't getfile from cloud", http.StatusInternalServerError)
			return
		}
		defer fl.Close()

		log.Info("file was gotten")

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
