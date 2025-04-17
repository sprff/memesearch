package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
)

func (a *API) CreateMeme(ctx context.Context, board models.BoardID, filename string, dsc map[string]string) (models.MemeID, error) {
	if err := a.aclPostMeme(ctx, board); err != nil {
		return "", fmt.Errorf("acl failed: %w", err)
	}

	meme := models.Meme{BoardID: board, Filename: filename, Description: dsc}
	id, err := a.storage.InsertMeme(ctx, meme)
	if err != nil {
		return "", fmt.Errorf("can't create meme: %w", err)
	}

	return id, nil
}

func (a *API) GetMemeByID(ctx context.Context, id models.MemeID) (models.Meme, error) {
	if err := a.aclGetMeme(ctx, id); err != nil {
		return models.Meme{}, fmt.Errorf("acl failed: %w", err)
	}

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

func (a *API) UpdateMeme(ctx context.Context, id models.MemeID, board *models.BoardID, filename *string, dsc *map[string]string) (models.Meme, error) {
	if err := a.aclUpdateMeme(ctx, id); err != nil {
		return models.Meme{}, fmt.Errorf("acl failed: %w", err)
	}

	meme, err := a.GetMemeByID(ctx, id)
	if err != nil {
		return models.Meme{}, fmt.Errorf("can't get meme: %w", err)
	}
	if dsc != nil {
		meme.Description = *dsc
	}
	if filename != nil {
		meme.Filename = *filename
	}
	if board != nil {
		meme.BoardID = *board
	}

	err = a.storage.UpdateMeme(ctx, meme)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMemeNotFound):
			return models.Meme{}, ErrMemeNotFound
		default:
			return models.Meme{}, fmt.Errorf("can't update meme: %w", err)
		}
	}

	meme, err = a.GetMemeByID(ctx, id)
	if err != nil {
		return models.Meme{}, fmt.Errorf("can't get meme: %w", err)
	}

	return meme, nil
}

func (a *API) DeleteMeme(ctx context.Context, id models.MemeID) error {
	if err := a.aclDeleteMeme(ctx, id); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}

	logger := slog.Default().With("from", "API.DeleteMeme")
	logger.InfoContext(ctx, "Started", "id", id)

	err := a.storage.DeleteMeme(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMemeNotFound):
			return ErrMemeNotFound
		default:
			return fmt.Errorf("can't delete meme: %w", err)
		}
	}
	return nil
}

func (a *API) ListMemes(ctx context.Context, offset, limit int, sortBy string) ([]models.Meme, error) {
	userID := GetUserID(ctx)
	if userID == "" {
		return nil, ErrUnauthorized
	}

	memes, err := a.storage.ListMemes(ctx, userID, offset, limit, sortBy)
	if err != nil {
		return nil, fmt.Errorf("can't list memes: %w", err)
	}

	return memes, nil
}
