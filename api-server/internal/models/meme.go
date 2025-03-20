package models

import (
	"context"
	"errors"
)

type MemeID string

type Meme struct {
	ID           MemeID            `json:"id"`
	BoardID      BoardID           `json:"board_id" `
	Filename     string            `json:"filename"`
	Descriptions map[string]string `json:"descriptions"`
}

type MemeRepo interface {
	InsertMeme(ctx context.Context, meme Meme) (MemeID, error)
	GetMemeByID(ctx context.Context, id MemeID) (Meme, error)
	GetMemesByBoardID(ctx context.Context, id BoardID, offset int, limit int) ([]Meme, error)
	UpdateMeme(ctx context.Context, meme Meme) error
	DeleteMeme(ctx context.Context, id MemeID) error
}

var ErrMemeNotFound = errors.New("Meme not found")
