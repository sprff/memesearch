package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"memesearch/internal/config"
	"memesearch/internal/models"
)

type MediaStore struct {
	client *YaClientS3
}

var _ models.MediaRepo = &MediaStore{}

func NewMediaStore(ctx context.Context, cfg config.S3Config) (*MediaStore, error) {
	c, err := GetClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &MediaStore{
		client: c,
	}, nil
}

func (s *MediaStore) GetMediaByID(ctx context.Context, id models.MediaID) (models.Media, error) {
	//TODO add ErrNoMediaFound
	media, err := s.client.GetObject(ctx, string(id))
	if err != nil {
		return models.Media{}, fmt.Errorf("can't get object: %w", err)
	}
	body, err := io.ReadAll(media)
	if err != nil {
		return models.Media{}, fmt.Errorf("can't read media: %w", err)
	}

	return models.Media{
		ID:   id,
		Body: body,
	}, nil
}

func (s *MediaStore) SetMediaByID(ctx context.Context, media models.Media) error {
	err := s.client.PutObject(ctx, string(media.ID), bytes.NewBuffer(media.Body))
	if err != nil {
		return fmt.Errorf("can't put media object: %w", err)
	}
	return nil
}
