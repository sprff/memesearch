package psql

import (
	"bytes"
	"context"
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"

	"github.com/jmoiron/sqlx"
)

var _ models.MediaRepo = &MediaStore{}

type MediaStore struct {
	db *sqlx.DB
}

func NewMediaStore(cfg config.DatabaseConfig) (MediaStore, error) {
	db, err := connect(cfg)
	if err != nil {
		return MediaStore{}, err
	}
	return MediaStore{db: db}, nil
}

type mediaBin struct {
	ID   models.MediaID `db:"id"`
	Body []byte         `db:"body"`
}

// GetMediaByID implements models.MediaRepo.
func (m *MediaStore) GetMediaByID(ctx context.Context, id models.MediaID) (models.Media, error) {

	med := mediaBin{}
	err := m.db.Get(&med, "SELECT * FROM medias WHERE id=$1", id)
	if err != nil {
		return models.Media{}, fmt.Errorf("can't select: %w", err)
	}
	return models.Media{ID: med.ID, Body: bytes.NewBuffer(med.Body)}, nil

}

// SetMediaByID implements models.MediaRepo.
func (m *MediaStore) SetMediaByID(ctx context.Context, media models.Media) error {
	med := mediaBin{ID: media.ID, Body: media.Body.Bytes()}
	_, err := m.db.Exec(`INSERT INTO medias (id, body) VALUES ($1, $2)
	ON CONFLICT (id) DO UPDATE SET  body=$2`, med.ID, med.Body)
	if err != nil {
		return fmt.Errorf("can't insert: %w", err)
	}

	return nil

}
