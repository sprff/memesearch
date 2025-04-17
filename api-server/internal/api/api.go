package api

import (
	"memesearch/internal/config"
	"memesearch/internal/searchranker"
	"memesearch/internal/storage"
)

type api struct {
	storage storage.Storage
	secrets config.SecretConfig
	ranker  searchranker.Ranker
}

func newApi(s storage.Storage, secrets config.SecretConfig, ranker searchranker.Ranker) *api {
	return &api{
		storage: s,
		secrets: secrets,
		ranker:  ranker,
	}
}
