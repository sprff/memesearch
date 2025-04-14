package main

import (
	"fmt"
	"log/slog"
	"memesearch/internal/api"
	"memesearch/internal/apiserver"
	"memesearch/internal/apiserver/middleware"
	"memesearch/internal/config"
	"memesearch/internal/contextlogger"
	"memesearch/internal/memesearcher"
	"memesearch/internal/storage"
	"net/http"
	"os"
	"time"
)

func main() {
	setLogger()
	cfg := getConfig()
	s, err := storage.New(cfg)
	processError("Failed to create storage", err)
	searcher, err := memesearcher.New(cfg)
	processError("Failed to create searcher", err)
	api := api.New(s, cfg.Secrets, searcher)
	server := apiserver.NewHandler(api, []middleware.Middleware{middleware.Logger(), middleware.Auth(api)})
	slog.Info("Run server", "port", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port), server)
	slog.Info("Server stopped", "err", err)
}

func setLogger() {
	logFolder := os.Getenv("LOG_FOLDER")
	date := time.Now().Format("2006-01-02")
	handler, err := contextlogger.NewContextHandler(fmt.Sprintf("%s/%s.log", logFolder, date), &slog.HandlerOptions{Level: slog.LevelDebug})
	processError("Can't get context handler", err)
	slog.SetDefault(slog.New(handler))
}

func getConfig() config.Config {
	configPath := os.Getenv("MS_API_CONFIG_PATH")
	if configPath == "" {
		slog.Error("Сonfig path is not specified")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(configPath)
	processError("Сan't load config", err)
	return cfg
}

func processError(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.String("error", err.Error()))
		os.Exit(1)
	}
}
