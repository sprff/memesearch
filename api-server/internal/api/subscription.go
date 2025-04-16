package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *API) Subscribe(ctx context.Context, sub models.Subsciption) error {
	err := a.storage.Subscribe(ctx, sub)
	if err != nil {
		return fmt.Errorf("can't subscribe: %w", err)
	}

	return nil
}

func (a *API) Unsubscribe(ctx context.Context, sub models.Subsciption) error {
	err := a.storage.Unsubscribe(ctx, sub)
	if err != nil {
		if err == models.ErrSubNotFound {
			return ErrSubNotFound
		}
		return fmt.Errorf("can't unsubscribe: %w", err)
	}
	return nil
}

