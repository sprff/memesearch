package psql

import (
	"context"
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"

	"github.com/jmoiron/sqlx"
)

var _ models.SubsciptionRepo = &SubStore{}

type SubStore struct {
	db *sqlx.DB
}

func NewSubStore(cfg config.DatabaseConfig) (*SubStore, error) {
	db, err := connect(cfg)
	if err != nil {
		return nil, err
	}
	return &SubStore{db: db}, nil
}

// Subscribe implements models.SubsciptionRepo.
func (s *SubStore) Subscribe(ctx context.Context,user models.UserID, board models.BoardID, role string) error {
	n := 0
	err := s.db.Get(&n, "SELECT COUNT(*) FROM subscriptions WHERE user_id=$1 AND board_id=$2", user, board)
	if err != nil {
		return fmt.Errorf("can't select: %w", err)
	}
	if n == 0 {
		_, err := s.db.Exec("INSERT INTO subscriptions (user_id, board_id, role) VALUES ($1, $2, $3)", user, board, role)
		if err != nil {
			return fmt.Errorf("can't insert: %w", err)
		}
		return nil
	}

	_, err = s.db.Exec("UPDATE subscriptions SET role=$3 WHERE user_id=$1 AND board_id=$2", user, board, role)
	if err != nil {
		return fmt.Errorf("can't insert: %w", err)
	}
	return nil
}

// Unsubscribe implements models.SubsciptionRepo.
func (s *SubStore) Unsubscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	res, err := s.db.Exec("DELETE FROM subscriptions WHERE user_id=$1 AND board_id=$2", user, board)
	if err != nil {
		return fmt.Errorf("can't delete: %w", err)
	}
	if err := zeroRows(res, models.ErrSubNotFound); err != nil {
		return err
	}
	return nil
}
