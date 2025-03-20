package models

import (
	"context"
	"errors"
)

type BoardID string

type Board struct {
	ID    BoardID `json:"id"    db:"id"`
	Owner UserID  `json:"owner" db:"owner_id"`
	Name  string  `json:"name"  db:"name"`
}

type BoardRepo interface {
	InsertBoard(ctx context.Context, board Board) (BoardID, error)
	GetBoardByID(ctx context.Context, id BoardID) (Board, error)
	UpdateBoard(ctx context.Context, board Board) error
	DeleteBoard(ctx context.Context, id BoardID) error
}

var ErrBoardNotFound = errors.New("Board not found")
