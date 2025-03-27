package models

import (
	"context"
)

type MediaID MemeID

type Media struct {
	ID   MediaID `json:"id"`
	Body []byte  `json:"body"`
}

type MediaRepo interface {
	GetMediaByID(ctx context.Context, id MediaID) (Media, error)
	SetMediaByID(ctx context.Context, media Media) error
}
