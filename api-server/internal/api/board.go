package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *api) CreateBoard(ctx context.Context, name string) (models.Board, error) {
	userID := GetUserID(ctx)
	board, err := a.storage.CreateBoard(ctx, userID, name)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't create board: %w", err)
	}

	return board, nil
}

func (a *api) GetBoardByID(ctx context.Context, id models.BoardID) (models.Board, error) {
	board, err := a.storage.GetBoardByID(ctx, id)
	if err != nil {
		if err == models.ErrBoardNotFound {
			return models.Board{}, ErrBoardNotFound
		}
		return models.Board{}, fmt.Errorf("can't get board: %w", err)
	}

	return board, nil
}

func (a *api) UpdateBoard(ctx context.Context, id models.BoardID, name *string, owner *models.UserID) (models.Board, error) {
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

func (a *api) DeleteBoard(ctx context.Context, id models.BoardID) (models.Board, error) {
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

func (a *api) ListBoards(ctx context.Context, offset, limit int, sortBy string) ([]models.Board, error) {
	userID := GetUserID(ctx)
	boards, err := a.storage.ListBoards(ctx, userID, offset, limit, sortBy)
	if err != nil {
		return nil, fmt.Errorf("can't list boards: %w", err)
	}
	return boards, nil
}
