package api

import "errors"

var (
	ErrUserNotFound  = errors.New("USER_NOT_FOUND")
	ErrBoardNotFound = errors.New("BOARD_NOT_FOUND")
	ErrMediaNotFound = errors.New("MEDIA_NOT_FOUND")
	ErrMemeNotFound  = errors.New("MEME_NOT_FOUND")
	ErrInvalidToken  = errors.New("INVALID_TOKEN")
	ErrForbidden     = errors.New("FORBIDDEN")
)
