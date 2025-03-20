package api

type ApiError interface {
	Status() string
	Error() string //Human readable
}

type ErrUserNotFound struct{}
type ErrBoardNotFound struct{}
type ErrMediaNotFound struct{}
type ErrMemeNotFound struct{}
type ErrUnexpectedError struct{ err error }

func (e ErrUserNotFound) Status() string    { return "USER_NOT_FOUND" }
func (e ErrBoardNotFound) Status() string   { return "BOARD_NOT_FOUND" }
func (e ErrMediaNotFound) Status() string   { return "MEDIA_NOT_FOUND" }
func (e ErrMemeNotFound) Status() string    { return "MEME_NOT_FOUND" }
func (e ErrUnexpectedError) Status() string { return "MEME_NOT_FOUND" }

func (e ErrUserNotFound) Error() string    { return "USER_NOT_FOUND" }
func (e ErrBoardNotFound) Error() string   { return "BOARD_NOT_FOUND" }
func (e ErrMediaNotFound) Error() string   { return "MEDIA_NOT_FOUND" }
func (e ErrMemeNotFound) Error() string    { return "MEME_NOT_FOUND" }
func (e ErrUnexpectedError) Error() string { return e.err.Error() }

func (e ErrUnexpectedError) Unwrap() error { return e.err }
