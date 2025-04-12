package apiserver

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
)

func (e Error) Error() string {
	return e.Message
}

func (e Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	return (t.Code == "" || t.Code == e.Code) &&
		(t.Message == "" || t.Message == e.Message) &&
		(t.Body == nil || (e.Body != nil && reflect.DeepEqual(*t.Body, *e.Body)))

}

var ErrMemeNotFound = Error{Code: "MEME_NOT_FOUND", Message: "Can't find meme"}
var ErrMediaNotFound = Error{Code: "MEDIA_NOT_FOUND", Message: "Can't find media"}
var ErrTooLarge = Error{Code: "FILE_TOO_LARGE", Message: "Request body too large"}
var ErrInvalidPagination = Error{Code: "INVALID_PAGINATION", Message: "Some of requrements doesn't meet page>=1;1<=pagesize<=100"}
var ErrUnsupportedMediaType = Error{Code: "UNSUPPORTED_MEDIA_TYPE", Message: "Unsupported media type"}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	slog.DebugContext(ctx, "Response error", "err", err)

	var invalidParam *InvalidParamFormatError
	switch {
	case errors.Is(err, ErrMediaNotFound),
		errors.Is(err, ErrMemeNotFound):

		w.WriteHeader(http.StatusNotFound)

	case errors.Is(err, ErrTooLarge),
		errors.Is(err, ErrInvalidPagination),
		errors.Is(err, ErrUnsupportedMediaType):

		w.WriteHeader(http.StatusBadRequest)

	case errors.As(err, &invalidParam):

		b := map[string]any{invalidParam.ParamName: invalidParam.Err.Error()}
		body, err := json.Marshal(Error{Code: "INVALID_REQUEST", Body: &b})
		if err != nil {
			slog.ErrorContext(ctx, "Can't marshall json", "err", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return

	default:
		slog.ErrorContext(ctx, "Unexpected error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := json.MarshalIndent(err, "", " ") //FIXME
	if err != nil {
		slog.ErrorContext(ctx, "Can't marshall json", "err", err)
	}
	w.Write(body)

}
