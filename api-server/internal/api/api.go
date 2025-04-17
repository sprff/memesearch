package api

import (
	"memesearch/internal/config"
	"memesearch/internal/memesearcher"
	"memesearch/internal/storage"
)

type api struct {
	storage  storage.Storage
	secrets  config.SecretConfig
	searcher memesearcher.Searcher
}

func newApi(s storage.Storage, secrets config.SecretConfig, searcher memesearcher.Searcher) *api {
	return &api{
		storage:  s,
		secrets:  secrets,
		searcher: searcher,
	}
}
