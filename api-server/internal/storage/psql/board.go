package psql

import (
	"context"
	"database/sql"
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"memesearch/internal/utils"

	"github.com/jmoiron/sqlx"
)

var _ models.BoardRepo = &BoardStore{}

type BoardStore struct {
	db *sqlx.DB
}

func NewBoardStore(cfg config.DatabaseConfig) (*BoardStore, error) {
	db, err := connect(cfg)
	if err != nil {
		return nil, err
	}
	return &BoardStore{db: db}, nil
}

// GetBoardByID implements models.BoardRepo.
func (b *BoardStore) GetBoardByID(ctx context.Context, id models.BoardID) (models.Board, error) {
	var board models.Board
	err := b.db.Get(&board, "SELECT * FROM boards WHERE id=$1", id)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return models.Board{}, models.ErrBoardNotFound
		default:
			return models.Board{}, fmt.Errorf("can't select: %w", err)
		}
	}
	return board, nil
}

// InsertBoard implements models.BoardRepo.
func (b *BoardStore) InsertBoard(ctx context.Context, board models.Board) (models.BoardID, error) {
	board.ID = models.BoardID(utils.GenereateUUIDv7())
	_, err := b.db.Exec("INSERT INTO boards (id, owner_id, name) VALUES ($1, $2, $3)", board.ID, board.Owner, board.Name)
	if err != nil {
		return "", fmt.Errorf("can't insert: %w", err)
	}
	return board.ID, nil

}

// UpdateBoard implements models.BoardRepo.
func (b *BoardStore) UpdateBoard(ctx context.Context, board models.Board) error {
	res, err := b.db.Exec("UPDATE boards SET owner_id = $2, name = $3 WHERE id=$1", board.ID, board.Owner, board.Name)
	if err != nil {
		return fmt.Errorf("can't update: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows: %w", err)
	}
	if rows == 0 {
		return models.ErrBoardNotFound
	}
	return nil
}

// DeleteBoard implements models.BoardRepo.
func (b *BoardStore) DeleteBoard(ctx context.Context, id models.BoardID) error {
	res, err := b.db.Exec("DELETE FROM boards WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("can't delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows: %w", err)
	}
	if rows == 0 {
		return models.ErrBoardNotFound
	}
	return nil
}
