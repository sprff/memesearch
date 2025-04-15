package models

import (
	"context"
)

type BoardID string

type Board struct {
	ID    BoardID `json:"id"    db:"id"`
	Owner UserID  `json:"owner" db:"owner_id"`
	Name  string  `json:"name"  db:"name"`
}

type BoardRepo interface {
	CreateBoard(ctx context.Context, owner UserID, name string) (Board, error)
	GetBoardByID(ctx context.Context, id BoardID) (Board, error)
	UpdateBoard(ctx context.Context, board Board) error
	DeleteBoard(ctx context.Context, id BoardID) error
	ListBoards(ctx context.Context, offset, limit int, sortBy string) ([]Board, error) 
}
