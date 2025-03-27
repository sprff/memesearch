package psql

import (
	"encoding/json"
	"fmt"
	"memesearch/internal/models"
	"time"
)

type psqlMeme struct {
	ID           models.MemeID  `db:"id"`
	BoardID      models.BoardID `db:"board_id"`
	Filename     string         `db:"filename"`
	Descriptions string         `db:"descriptions"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

func convertModelsMeme(m models.Meme) (psqlMeme, error) {
	data, err := json.Marshal(m.Descriptions)
	if err != nil {
		return psqlMeme{}, fmt.Errorf("can't marshal: %w", err)
	}
	return psqlMeme{
		ID:           m.ID,
		BoardID:      m.BoardID,
		Filename:     m.Filename,
		Descriptions: string(data),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}

func convertPsqlMeme(m psqlMeme) (models.Meme, error) {
	data := map[string]string{}
	err := json.Unmarshal([]byte(m.Descriptions), &data)
	if err != nil {
		return models.Meme{}, fmt.Errorf("can't unmarshal: %w", err)
	}
	return models.Meme{
		ID:           m.ID,
		BoardID:      m.BoardID,
		Filename:     m.Filename,
		Descriptions: data,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}
