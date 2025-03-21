package httpserver

import (
	"fmt"
	"memesearch/internal/api"
	"memesearch/internal/config"
	"memesearch/internal/storage"
	"net/http"
)

func New(cfg config.Config) (http.Handler, error) {
	storage, err := storage.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create storage: %w", err)
	}
	api := api.New(storage, cfg.Secrets)
	router := GetRouter(api)
	return router, nil
}
