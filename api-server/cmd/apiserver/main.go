package main

import (
	"fmt"
	"log/slog"
	"memesearch/internal/config"
	"memesearch/internal/contextlogger"
	"memesearch/internal/httpserver"
	"net/http"
	"os"
	"time"
)

func main() {
	setLogger()
	cfg := getConfig()
	server, err := httpserver.New(cfg)
	processError("Failed to create server", err)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port), server)
	slog.Info("Server stopped")
}

func setLogger() {
	handler, err := contextlogger.NewContextHandler(fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02")), &slog.HandlerOptions{Level: slog.LevelDebug})
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
