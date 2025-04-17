package api

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound  = errors.New("USER_NOT_FOUND")
	ErrBoardNotFound = errors.New("BOARD_NOT_FOUND")
	ErrMediaNotFound = errors.New("MEDIA_NOT_FOUND")
	ErrMemeNotFound  = errors.New("MEME_NOT_FOUND")
	ErrSubNotFound   = errors.New("SUB_NOT_FOUND")

	ErrInvalidToken = errors.New("INVALID_TOKEN")
	ErrForbidden    = errors.New("FORBIDDEN")
	ErrUnauthorized = errors.New("UNAUTHORIZED")
	ErrLoginExists  = errors.New("LOGIN_EXISTS")
)

type ErrInvalid struct {
	Param  string
	Reason string
}

func (e ErrInvalid) Error() string { return fmt.Sprintf("invalid %s: %s", e.Param, e.Reason) }
func (e ErrInvalid) Is(target error) bool {
	t, ok := target.(ErrInvalid)
	if !ok {
		return false
	}
	return (t.Reason == "" || t.Reason == e.Reason) &&
		(t.Param == "" || t.Param == e.Param)
}

