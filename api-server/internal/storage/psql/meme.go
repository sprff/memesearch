package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"memesearch/internal/utils"

	"github.com/jmoiron/sqlx"
)

var _ models.MemeRepo = &MemeStore{}

type MemeStore struct {
	db *sqlx.DB
}

func NewMemeStore(cfg config.DatabaseConfig) (*MemeStore, error) {
	db, err := connect(cfg)
	if err != nil {
		return nil, err
	}
	return &MemeStore{db: db}, nil
}

// InsertMeme implements models.MemeRepo.
func (m *MemeStore) InsertMeme(ctx context.Context, meme models.Meme) (models.MemeID, error) {
	meme.ID = models.MemeID(utils.GenereateUUIDv7())
	mp, err := convertModelsMeme(meme)
	if err != nil {
		return "", fmt.Errorf("can't convert: %w", err)
	}
	_, err = m.db.Exec("INSERT INTO memes (id, board_id, descriptions, filename, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", mp.ID, mp.BoardID, mp.Descriptions, mp.Filename)
	if err != nil {
		return "", fmt.Errorf("can't insert: %w", err)
	}

	return meme.ID, nil
}

// GetMemeByID implements models.MemeRepo.
func (m *MemeStore) GetMemeByID(ctx context.Context, id models.MemeID) (models.Meme, error) {
	var mp psqlMeme
	err := m.db.Get(&mp, "SELECT * FROM memes WHERE id=$1", id)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return models.Meme{}, models.ErrMemeNotFound
		default:
			return models.Meme{}, fmt.Errorf("can't select: %w", err)
		}
	}

	meme, err := convertPsqlMeme(mp)
	if err != nil {
		return models.Meme{}, fmt.Errorf("can't convert: %w", err)
	}
	return meme, nil
}

// GetMemesByBoardID implements models.MemeRepo.
func (m *MemeStore) GetMemesByBoardID(ctx context.Context, id models.BoardID, offset int, limit int) ([]models.Meme, error) {
	var mps []psqlMeme
	err := m.db.Select(&mps, "SELECT * FROM memes WHERE board_id=$1 ORDER BY id OFFSET $2 LIMIT $3", id, offset, limit)
	if err != nil {
		return []models.Meme{}, fmt.Errorf("can't select: %w", err)
	}
	slog.DebugContext(ctx, "select result", "mps", mps)
	memes := make([]models.Meme, 0, len(mps))
	for _, mp := range mps {
		meme, err := convertPsqlMeme(mp)

		if err != nil {
			return []models.Meme{}, fmt.Errorf("can't convert: %w", err)
		}
		memes = append(memes, meme)
	}
	slog.DebugContext(ctx, "convert result", "memes", mps)

	return memes, nil
}

// UpdateMeme implements models.MemeRepo.
func (m *MemeStore) UpdateMeme(ctx context.Context, meme models.Meme) error {
	mp, err := convertModelsMeme(meme)
	if err != nil {
		return fmt.Errorf("can't convert: %w", err)
	}
	res, err := m.db.Exec("UPDATE memes SET board_id = $2, descriptions = $3, filename = $4, updated_at=CURRENT_TIMESTAMP WHERE id=$1", mp.ID, mp.BoardID, mp.Descriptions, mp.Filename)
	if err != nil {
		return fmt.Errorf("can't update: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows: %w", err)
	}
	if rows == 0 {
		return models.ErrMemeNotFound
	}
	return nil
}

func (m *MemeStore) DeleteMeme(ctx context.Context, id models.MemeID) error {
	res, err := m.db.Exec("DELETE FROM memes WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("can't delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows: %w", err)
	}
	if rows == 0 {
		return models.ErrMemeNotFound
	}
	return nil
}

func (m *MemeStore) ListMemes(ctx context.Context, request models.ListMemesRequest) ([]models.Meme, error) {
	var mps []psqlMeme
	err := m.db.Select(&mps, "SELECT * FROM memes ORDER BY id OFFSET $1 LIMIT $2", request.Offset, request.Limit)
	if err != nil {
		return []models.Meme{}, fmt.Errorf("can't select: %w", err)
	}
	slog.DebugContext(ctx, "select result", "mps", mps)
	memes := make([]models.Meme, 0, len(mps))
	for _, mp := range mps {
		meme, err := convertPsqlMeme(mp)
		if err != nil {
			return []models.Meme{}, fmt.Errorf("can't convert: %w", err)
		}
		memes = append(memes, meme)
	}
	slog.DebugContext(ctx, "convert result", "memes", mps)

	return memes, nil
}
