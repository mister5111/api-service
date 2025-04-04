package main

import (
	"api-service/src/config"
	"api-service/src/handlers/delete"
	"api-service/src/handlers/example"
	"api-service/src/handlers/save"
	"api-service/src/handlers/show"
	"api-service/src/lib/slogpretty"
	"api-service/src/storage/sqlite"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.ConfigLoad()

	log := setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Storage initialization error", slog.String("err", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", example.Example(log))
	router.Post("/save", save.New(log, storage))
	router.Post("/del", delete.Del(log, storage))
	router.Get("/show", show.Show(log, storage))
	router.Get("/all", show.ShowAll(log, storage))

	log.Info("Starting server", slog.String("Address", cfg.Address), slog.String("LogLevel", cfg.Env))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Server error", slog.String("err", err.Error()))
		os.Exit(1)
	}

	log.Error("Server down")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slogpretty.SetupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
