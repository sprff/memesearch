package models

import (
	"context"
	"time"
)

type MemeID string

type Meme struct {
	ID          MemeID            `json:"id"`
	BoardID     BoardID           `json:"board_id"`
	Filename    string            `json:"filename"`
	Description map[string]string `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type MemeRepo interface {
	InsertMeme(ctx context.Context, meme Meme) (MemeID, error)
	GetMemeByID(ctx context.Context, id MemeID) (Meme, error)
	GetMemesByBoardID(ctx context.Context, id BoardID, offset int, limit int) ([]Meme, error)
	UpdateMeme(ctx context.Context, meme Meme) error
	DeleteMeme(ctx context.Context, id MemeID) error
}
