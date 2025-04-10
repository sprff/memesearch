package storage

import (
	"context"
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"memesearch/internal/storage/psql"
	"memesearch/internal/storage/s3"
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
	s.MediaRepo, err = s3.NewMediaStore(context.TODO(), cfg.S3)
	if err != nil {
		return Storage{}, fmt.Errorf("can't load media store: %w", err)
	}
	return s, nil
}
