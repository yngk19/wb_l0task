package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yngk19/wb_l0task/internal/config"
	"github.com/yngk19/wb_l0task/internal/repository"
	"github.com/yngk19/wb_l0task/internal/net/nats"
	get_time "github.com/yngk19/wb_l0task/internal/net/http/time"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Service.Env)

	log.Info("Orders service is starting!", slog.String("env", cfg.Service.Env))

	_, err := repository.NewDB(cfg.DB, log)
	if err != nil {
		log.Error("Failed connection to storage!: %s", err)
		os.Exit(1)
	}
		
	err = nats.Connect(cfg.Nats, log)
	if err != nil {
		log.Error("NATS", slog.String("Error", err))
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/time", get_time.Time(log))

	log.Info("Starting the http server on", slog.Int("port", cfg.Service.HTTPServer.Port))

	srv := &http.Server{
		Addr:         cfg.Service.HTTPServer.Address,
		Handler:      r,
		IdleTimeout:  cfg.Service.HTTPServer.IddleTimeout,
		ReadTimeout:  cfg.Service.HTTPServer.Timeout,
		WriteTimeout: cfg.Service.HTTPServer.Timeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start http server!")
	}

	log.Error("Server stoped!")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
