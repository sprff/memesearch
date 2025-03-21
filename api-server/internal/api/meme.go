package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
)

func (a *API) CreateMeme(ctx context.Context, meme models.Meme) (models.MemeID, error) {
	logger := slog.Default().With("from", "API.CreateMeme")
	logger.InfoContext(ctx, "Started")

	id, err := a.storage.InsertMeme(ctx, meme)
	if err != nil {
		return "", fmt.Errorf("can't create meme: %w", err)
	}
	logger.DebugContext(ctx, "Meme inserted", "id", id)
	return id, nil
}

func (a *API) GetMemeByID(ctx context.Context, id models.MemeID) (models.Meme, error) {
	logger := slog.Default().With("from", "API.GetMeme")
	logger.InfoContext(ctx, "Started", "id", id)

	meme, err := a.storage.GetMemeByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMemeNotFound):
			return models.Meme{}, ErrMemeNotFound
		default:
			return models.Meme{}, fmt.Errorf("can't get meme: %w", err)
		}
	}
	return meme, nil
}

func (a *API) UpdateMeme(ctx context.Context, meme models.Meme) error {
	logger := slog.Default().With("from", "API.UpdateMeme")
	logger.InfoContext(ctx, "Started", "id", meme.ID)

	err := a.storage.UpdateMeme(ctx, meme)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMemeNotFound):
			return ErrMemeNotFound
		default:
			return fmt.Errorf("can't update meme: %w", err)
		}
	}
	return nil
}

func (a *API) DeleteMeme(ctx context.Context, id models.MemeID) error {
	logger := slog.Default().With("from", "API.DeleteMeme")
	logger.InfoContext(ctx, "Started", "id", id)

	err := a.storage.DeleteMeme(ctx, id)
	switch {
	case errors.Is(err, models.ErrMemeNotFound):
		return ErrMemeNotFound
	default:
		return fmt.Errorf("can't delete meme: %w", err)
	}
}
