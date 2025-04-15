package apiserver

import (
	"encoding/json"
	"errors"
	"log/slog"
	"memesearch/internal/api"
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
	return (t.Id == "" || t.Id == e.Id) &&
		(t.Message == "" || t.Message == e.Message) &&
		(t.Body == nil || (e.Body != nil && reflect.DeepEqual(*t.Body, *e.Body)))

}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	slog.DebugContext(ctx, "Response error", "err", err)

	var resultErr Error
	var invalidParam *InvalidParamFormatError

	switch {
	case errors.Is(err, api.ErrMediaNotFound),
		errors.Is(err, api.ErrMemeNotFound),
		errors.Is(err, api.ErrUserNotFound),
		errors.Is(err, api.ErrBoardNotFound):

		w.WriteHeader(http.StatusNotFound)
		resultErr = Error{Id: err.Error()}

	case errors.As(err, &invalidParam):
		b := map[string]any{invalidParam.ParamName: invalidParam.Err.Error()}
		resultErr = Error{Id: "INVALID_REQUEST", Body: &b}
		w.WriteHeader(http.StatusBadRequest)

	default:
		slog.ErrorContext(ctx, "Unexpected error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.MarshalIndent(resultErr, "", " ")
	if err != nil {
		slog.ErrorContext(ctx, "Can't marshall json", "err", err)
	}
	w.Write(body)

}
