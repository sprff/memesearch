package api

import (
	"memesearch/internal/config"
	"memesearch/internal/storage"
)

type API struct {
	storage storage.Storage
	secrets config.SecretConfig
}

func New(s storage.Storage, secrets config.SecretConfig) *API {
	return &API{
		storage: s,
		secrets: secrets,
	}
}
