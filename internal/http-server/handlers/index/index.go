package index

import (
	grpcclient "client/internal/app/grpcClient"
	"context"
	"html/template"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, client *grpcclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := client.GetFullData(context.Background())
		if err != nil {
			http.Error(w, "Can't get data", http.StatusInternalServerError)
			return
		}

		tpl := template.Must(template.ParseFiles("../../internal/web-site/public/index.html"))
		tpl.Execute(w, data)
	}
}
