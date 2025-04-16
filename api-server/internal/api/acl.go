package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *API) aclGetBoard(ctx context.Context, id models.BoardID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}

	return nil
}

func (a *API) aclUpdateBoard(ctx context.Context, id models.BoardID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	board, err := a.GetBoardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}
	if board.Owner != userID {
		return ErrForbidden
	}
	return nil
}

func (a *API) aclDeleteBoard(ctx context.Context, id models.BoardID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	board, err := a.GetBoardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}
	if board.Owner != userID {
		return ErrForbidden
	}

	return nil
}
