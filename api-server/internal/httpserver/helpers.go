package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	apiservice "memesearch/internal/api"
	"memesearch/internal/contextlogger"
	"memesearch/internal/utils"
	"net/http"

	"github.com/go-chi/render"
)

type handlerWithError = func(w http.ResponseWriter, r *http.Request, ctx context.Context, a *apiservice.API) (any, error)

func handlerWrapper(handle handlerWithError, a *apiservice.API) http.HandlerFunc {
	type httpAnswer struct {
		Status  string `json:"status"`
		Data    any    `json:"data,omitempty"`
		ErrData any    `json:"err_data,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		ctx = contextlogger.AppendCtx(ctx, slog.String("request_id", utils.GenereateUUIDv7()))
		defer cancel()

		res, err := handle(w, r, ctx, a)
		if err != nil {
			status := "UNEXPECTED_ERROR"
			if aerr := UnwrapApiError(err); aerr != nil {
				status = aerr.Status()
			}
			slog.ErrorContext(ctx, "Error in handle", "error", err.Error())
			render.JSON(w, r, httpAnswer{
				Status:  status,
				ErrData: err,
			})
			return
		}
		render.JSON(w, r, httpAnswer{
			Status: "OK",
			Data:   res,
		})

	}
}

func UnwrapApiError(err error) apiservice.ApiError {
	type unwrapable interface {
		Unwrap() error
	}
	for {
		uerr, ok := err.(unwrapable)
		if !ok {
			return nil
		}
		aerr, ok := uerr.(apiservice.ApiError)
		if ok {
			return aerr
		}
		err = uerr.Unwrap()
	}
}

func readBody[T any](r *http.Request, input *T) error {
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		return fmt.Errorf("can't read body: %w", err)
	}

	err = json.Unmarshal(body, input)
	if err != nil {
		// var jerr *json.UnmarshalTypeError
		// if errors.As(err, &jerr) {
		// 	return nil, ErrBadParams{fmt.Sprintf("%s expected to be %v", jerr.Field, jerr.Type)}
		// }
		// return nil, ErrBadParams{"body"}
		return fmt.Errorf("can't unmarshal body: %w", err)
	}
	return nil
}
