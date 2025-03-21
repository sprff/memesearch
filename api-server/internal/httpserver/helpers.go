package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func readBody[T any](r *http.Request, input *T) error {
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		return fmt.Errorf("can't read body: %w", err)
	}

	err = json.Unmarshal(body, input)
	if err != nil {
		var jerr *json.UnmarshalTypeError
		if errors.As(err, &jerr) {
			return fmt.Errorf("can't unmarshal body: %w", ErrInvalidInput{fmt.Sprintf("%s expected to be %v", jerr.Field, jerr.Type)})
		}
		return fmt.Errorf("can't unmarshal body: %w", ErrInvalidInput{"body"})
	}
	return nil
}

type httpAnswer struct {
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
	ErrData any    `json:"err_data,omitempty"`
}

func renderError(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		slog.ErrorContext(r.Context(), err.Error())
		err := parseError(err)
		w.WriteHeader(err.RespCode)
		tmp := httpAnswer{
			Status:  err.Status,
			ErrData: err.Unwrap(),
		}
		render.JSON(w, r, tmp)

		return true
	}
	return false
}

func renderOK(w http.ResponseWriter, r *http.Request, res any) {
	render.JSON(w, r, httpAnswer{
		Status: "OK",
		Data:   res,
	})
}
