package api

import (
	"memesearch/internal/config"
	"memesearch/internal/memesearcher"
	"memesearch/internal/storage"
)

type API struct {
	storage  storage.Storage
	secrets  config.SecretConfig
	searcher memesearcher.Searcher
}

func New(s storage.Storage, secrets config.SecretConfig, searcher memesearcher.Searcher) *API {
	return &API{
		storage:  s,
		secrets:  secrets,
		searcher: searcher,
	}
}
