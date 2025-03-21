package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
)

func (a *API) GetMedia(ctx context.Context, id models.MediaID) (models.Media, ApiError) {
	logger := slog.Default().With("from", "API.GetMedia")
	logger.InfoContext(ctx, "Started", "id", id)

	media, err := a.storage.GetMediaByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMediaNotFound):
			return models.Media{}, ErrMediaNotFound{}
		default:
			return models.Media{}, ErrUnexpectedError{fmt.Errorf("can't get media: %w", err)}
		}
	}

	return media, nil
}

func (a *API) SetMedia(ctx context.Context, media models.Media) ApiError {
	logger := slog.Default().With("from", "API.SetMedia")
	logger.InfoContext(ctx, "Started", "id", media.ID)

	err := a.storage.SetMediaByID(ctx, media)
	if err != nil {
		return ErrUnexpectedError{fmt.Errorf("can't set media: %w", err)}
	}

	return nil
}
