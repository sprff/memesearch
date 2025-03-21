package storage

import (
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"memesearch/internal/storage/psql"
)

type Storage struct {
	models.BoardRepo
	models.MemeRepo
	models.MediaRepo
	models.UserRepo
}

func New(cfg config.Config) (s Storage, err error) {
	s.BoardRepo, err = psql.NewBoardStore(cfg.Database)
	if err != nil {
		return Storage{}, fmt.Errorf("can't load board store: %w", err)
	}
	s.UserRepo, err = psql.NewUserStore(cfg.Database)
	if err != nil {
		return Storage{}, fmt.Errorf("can't load user store: %w", err)
	}
	s.MemeRepo, err = psql.NewMemeStore(cfg.Database)
	if err != nil {
		return Storage{}, fmt.Errorf("can't load meme store: %w", err)
	}
	s.MediaRepo, err = psql.NewMediaStore(cfg.Database)
	if err != nil {
		return Storage{}, fmt.Errorf("can't load media store: %w", err)
	}
	return s, nil
}
