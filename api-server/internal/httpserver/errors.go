package httpserver

import (
	"errors"
	"fmt"
	"net/http"

	apiservice "memesearch/internal/api"
)

type httpError struct {
	RespCode int
	Status   string

	err error
}

func (h httpError) Error() string {
	return h.err.Error()
}

func (h httpError) Unwrap() error {
	return h.err
}

type ErrInvalidInput struct {
	Reason string `json:"reason"`
}

func (e ErrInvalidInput) Error() string { return fmt.Sprintf("invalid input: %s", e.Reason) }
func (e ErrInvalidInput) Is(taget error) bool {
	t, ok := taget.(ErrInvalidInput)
	if !ok {
		return false
	}
	return t.Reason == "" || t.Reason == e.Reason
}

func parseError(err error) httpError {
	h := httpError{err: err}
	var (
		ii ErrInvalidInput
	)
	switch {
	case errors.Is(err, apiservice.ErrBoardNotFound):
		h.RespCode = http.StatusNotFound
		h.Status = "BOARD_NOT_FOUND"
	case errors.Is(err, apiservice.ErrMediaNotFound):
		h.RespCode = http.StatusNotFound
		h.Status = "MEDIA_NOT_FOUND"
	case errors.Is(err, apiservice.ErrMemeNotFound):
		h.RespCode = http.StatusNotFound
		h.Status = "MEME_NOT_FOUND"
	case errors.Is(err, apiservice.ErrUserNotFound):
		h.RespCode = http.StatusNotFound
		h.Status = "USER_NOT_FOUND"
	case errors.Is(err, ErrMediaIsRequired):
		h.RespCode = http.StatusNotFound
		h.Status = "MEDIA_IS_REQUIRED"
	case errors.As(err, &ii):
		h.RespCode = http.StatusBadRequest
		h.Status = "INVALID_INPUT"
		h.err = ii
	default:
		h.RespCode = http.StatusInternalServerError
		h.Status = "UNEXPECTED_ERROR"
	}
	return h
}
