package httpserver

import (
	grpcclient "client/internal/app/grpcClient"
	"client/internal/http-server/handlers/index"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTP_server struct {
	serv *http.Server
	log  *slog.Logger
}

func NewServer(log *slog.Logger, address string, timeout, idle_timeout time.Duration, client *grpcclient.Client) *HTTP_server {
	router := setupRouter(log, client)

	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  idle_timeout,
	}

	return &HTTP_server{
		serv: srv,
		log:  log,
	}
}

func setupRouter(log *slog.Logger, client *grpcclient.Client) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("../../internal/web-site/public/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", index.New(log, client))

	return router
}

func (h *HTTP_server) RunServer() {
	h.log.Info("starting http-server")
	if err := h.serv.ListenAndServe(); err != nil {
		h.log.Error("failed in ListenAndServe")
	}
}

//TODO: graceful stop
