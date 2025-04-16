package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *API) CreateBoard(ctx context.Context, name string) (models.Board, error) {
	userID := GetUserID(ctx)
	if userID == "" {
		return models.Board{}, ErrUnauthorized
	}

	board, err := a.storage.CreateBoard(ctx, userID, name)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't create board: %w", err)
	}
	return board, nil
}

func (a *API) GetBoardByID(ctx context.Context, id models.BoardID) (models.Board, error) {
	if err := a.aclGetBoard(ctx, id); err != nil {
		return models.Board{}, fmt.Errorf("acl failed: %w", err)
	}

	board, err := a.storage.GetBoardByID(ctx, id)
	if err != nil {
		if err == models.ErrBoardNotFound {
			return models.Board{}, ErrBoardNotFound
		}
		return models.Board{}, fmt.Errorf("can't get board: %w", err)
	}

	return board, nil
}

func (a *API) UpdateBoard(ctx context.Context, id models.BoardID, name *string, owner *models.UserID) (models.Board, error) {
	if err := a.aclUpdateBoard(ctx, id); err != nil {
		return models.Board{}, fmt.Errorf("acl failed: %w", err)
	}

	board, err := a.GetBoardByID(ctx, id)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't get init board: %w", err)
	}

	if name != nil {
		board.Name = *name
	}
	if owner != nil {
		board.Owner = *owner
	}

	err = a.storage.UpdateBoard(ctx, board)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't update board: %w", err)
	}

	board, err = a.storage.GetBoardByID(ctx, id)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't get board: %w", err)
	}

	return board, nil
}

func (a *API) DeleteBoard(ctx context.Context, id models.BoardID) (models.Board, error) {
	if err := a.aclDeleteBoard(ctx, id); err != nil {
		return models.Board{}, fmt.Errorf("acl failed: %w", err)
	}

	board, err := a.storage.GetBoardByID(ctx, id)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't get board: %w", err)
	}
	err = a.storage.DeleteBoard(ctx, id)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't delete board: %w", err)
	}
	return board, nil
}

func (a *API) ListBoards(ctx context.Context, offset, limit int, sortBy string) ([]models.Board, error) {
	userID := GetUserID(ctx)
	boards, err := a.storage.ListBoards(ctx, userID, offset, limit, sortBy)
	if err != nil {
		return nil, fmt.Errorf("can't list boards: %w", err)
	}
	return boards, nil
}
