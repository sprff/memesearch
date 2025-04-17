package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
)

func (a *api) GetMedia(ctx context.Context, id models.MediaID) (models.Media, error) {
	logger := slog.Default().With("from", "api.GetMedia")
	logger.InfoContext(ctx, "Started", "id", id)

	media, err := a.storage.GetMediaByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMediaNotFound):
			return models.Media{}, ErrMediaNotFound
		default:
			return models.Media{}, fmt.Errorf("can't get media: %w", err)
		}
	}

	return media, nil
}

func (a *api) SetMedia(ctx context.Context, media models.Media) error {
	logger := slog.Default().With("from", "api.SetMedia")
	logger.InfoContext(ctx, "Started", "id", media.ID)

	err := a.storage.SetMediaByID(ctx, media)
	if err != nil {
		return fmt.Errorf("can't set media: %w", err)
	}

	return nil
}
