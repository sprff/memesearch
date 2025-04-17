package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *API) Subscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	if err := a.aclSubscribe(ctx, user, board, role); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}

	err := a.storage.Subscribe(ctx, user, board, role)
	if err != nil {
		return fmt.Errorf("can't subscribe: %w", err)
	}

	return nil
}

func (a *API) Unsubscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	if err := a.aclUnsubscribe(ctx, user, board, role); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}

	err := a.storage.Unsubscribe(ctx, user, board, role)
	if err != nil {
		if err == models.ErrSubNotFound {
			return ErrSubNotFound
		}
		return fmt.Errorf("can't unsubscribe: %w", err)
	}
	return nil
}
