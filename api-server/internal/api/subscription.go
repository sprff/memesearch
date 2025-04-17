package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *api) Subscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	err := a.storage.Subscribe(ctx, user, board, role)
	if err != nil {
		return fmt.Errorf("can't subscribe: %w", err)
	}

	return nil
}

func (a *api) Unsubscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	err := a.storage.Unsubscribe(ctx, user, board, role)
	if err != nil {
		if err == models.ErrSubNotFound {
			return ErrSubNotFound
		}
		return fmt.Errorf("can't unsubscribe: %w", err)
	}
	return nil
}
