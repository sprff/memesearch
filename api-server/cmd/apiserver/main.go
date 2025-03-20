package main

import (
	"fmt"
	"log/slog"
	"memesearch/internal/api"
	"memesearch/internal/config"
	"memesearch/internal/contextlogger"
	"memesearch/internal/httpserver"
	"memesearch/internal/storage"
	"memesearch/internal/storage/psql"
	"net/http"
	"os"
	"time"
)

func main() {
	setLogger()

	cfg := getConfig()
	storage := getStorage(cfg)
	api := api.New(storage, cfg.Secrets)
	router := httpserver.GetRouter(api, cfg.Server)

	slog.Info("Server started",
		slog.Int("Port", cfg.Server.Port),
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.Timeout,
	}

	err := srv.ListenAndServe()
	processError("Failed to start server", err)
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

func getStorage(cfg config.Config) storage.Storage {
	board, err := psql.NewBoardStore(cfg.Database)
	processError("can't load board store", err)
	user, err := psql.NewUserStore(cfg.Database)
	processError("can't load user store", err)
	meme, err := psql.NewMemeStore(cfg.Database)
	processError("can't load meme store", err)
	media, err := psql.NewMediaStore(cfg.Database)
	processError("can't load media store", err)

	return storage.Storage{
		BoardRepo: &board,
		UserRepo:  &user,
		MemeRepo:  &meme,
		MediaRepo: &media,
	}
}

func processError(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.String("error", err.Error()))
		os.Exit(1)
	}
}
