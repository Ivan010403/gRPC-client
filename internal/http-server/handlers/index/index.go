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
		client.GetFullData(context.Background())
		tpl := template.Must(template.ParseFiles("../../internal/web-site/public/index.html"))
		tpl.Execute(w, nil)
	}
}
