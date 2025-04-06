package models

import (
	"errors"
	"fmt"
)

// Board
var ErrBoardNotFound = errors.New("Board not found")

// Media
var ErrMediaNotFound = errors.New("Media not found")
var ErrMediaIsRequired = errors.New("Media is required")

// Meme
var ErrMemeNotFound = errors.New("Meme not found")

// User
var ErrUserNotFound = errors.New("User not found")
var ErrUserLoginAlreadyExists = errors.New("User with this login already exists")

// Api
type ErrInvalidInput struct {
	Reason string `json:"reason"`
}

// Impls
func (e ErrInvalidInput) Error() string { return fmt.Sprintf("invalid input: %s", e.Reason) }
func (e ErrInvalidInput) Is(taget error) bool {
	t, ok := taget.(ErrInvalidInput)
	if !ok {
		return false
	}
	return t.Reason == "" || t.Reason == e.Reason
}
