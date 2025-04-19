package models

import (
	"time"
)

type MemeID string

type Meme struct {
	ID           MemeID            `json:"id"`
	BoardID      BoardID           `json:"board_id"`
	Filename     string            `json:"filename"`
	Descriptions map[string]string `json:"descriptions"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type ScoredMeme struct {
	Score float64 `json:"score"`
	Meme  Meme    `json:"meme"`
}
