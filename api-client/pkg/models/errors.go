package models

import (
	"errors"
	"fmt"
)

// Board
var (
	ErrBoardNotFound   = errors.New("Board not found")
	ErrSubNotFound     = errors.New("Sub not found")
	ErrMediaNotFound   = errors.New("Media not found")
	ErrMediaIsRequired = errors.New("Media is required")
	ErrMemeNotFound    = errors.New("Meme not found")
	ErrUserNotFound    = errors.New("User not found")
	ErrLoginExists     = errors.New("User with this login already exists")
	ErrUnauthorized    = errors.New("Unauthorized")
	ErrForbidden       = errors.New("Forbidden")
)

// Api
type ErrInvalidInput struct {
	Param  string
	Reason string
}

// Impls
func (e ErrInvalidInput) Error() string { return fmt.Sprintf("invalid %s: %s", e.Param, e.Reason) }
func (e ErrInvalidInput) Is(taget error) bool {
	t, ok := taget.(ErrInvalidInput)
	if !ok {
		return false
	}
	return (t.Reason == "" || t.Reason == e.Reason) &&
		(t.Param == "" || t.Param == e.Param)

}
