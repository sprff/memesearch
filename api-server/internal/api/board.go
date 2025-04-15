package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *API) CreateBoard(ctx context.Context, owner models.UserID, name string) (models.Board, error) {
	board, err := a.storage.CreateBoard(ctx, owner, name)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't create board: %w", err)
	}
	return board, nil
}

func (a *API) GetBoardByID(ctx context.Context, id models.BoardID) (models.Board, error) {
	board, err := a.storage.GetBoardByID(ctx, id)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't get board: %w", err)
	}
	return board, nil
}

func (a *API) UpdateBoard(ctx context.Context, board models.Board) (models.Board, error) {
	err := a.storage.UpdateBoard(ctx, board)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't update board: %w", err)
	}
	board, err = a.storage.GetBoardByID(ctx, board.ID)
	if err != nil {
		return models.Board{}, fmt.Errorf("can't get board: %w", err)
	}

	return board, nil
}
func (a *API) DeleteBoard(ctx context.Context, id models.BoardID) (models.Board, error) {
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
	boards, err := a.storage.ListBoards(ctx, offset, limit, sortBy)
	if err != nil {
		return nil, fmt.Errorf("can't list boards: %w", err)
	}
	return boards, nil
}
