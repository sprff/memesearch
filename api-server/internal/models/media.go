package models

import (
	"bytes"
	"context"
	"errors"
)

type MediaID MemeID

type Media struct {
	ID   MediaID `json:"id"`
	Body *bytes.Buffer
}

type MediaRepo interface {
	GetMediaByID(ctx context.Context, id MediaID) (Media, error)
	SetMediaByID(ctx context.Context, media Media) error
}

// Errors

var ErrMediaNotFound = errors.New("Media not found")
