package httpserver

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	apiservice "memesearch/internal/api"
	"memesearch/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetMedia(a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := slog.Default().With("from", "Server.PutMedia")
		id := models.MediaID(chi.URLParam(r, "id"))
		logger.InfoContext(ctx, "Started", "id", id)

		media, err := a.GetMedia(ctx, models.MediaID(id))
		if err != nil {
			renderError(w, r, fmt.Errorf("can't get media: %w", err))
			return
		}

		if _, err := io.Copy(w, bytes.NewBuffer(media.Body)); err != nil {
			renderError(w, r, fmt.Errorf("can't copy media: %w", err))
			return
		}
	}
}

func PutMedia(a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := slog.Default().With("from", "Server.PutMedia")
		id := models.MediaID(chi.URLParam(r, "id"))
		logger.InfoContext(ctx, "Started", "id", id)

		file, _, err := r.FormFile("media")
		if err != nil {
			renderError(w, r, ErrMediaIsRequired)
			return
		}
		body, err := io.ReadAll(file)
		if err != nil {
			renderError(w, r, ErrInvalidInput{"can't read media file"})
			return

		}
		err = a.SetMedia(ctx, models.Media{ID: id, Body: body})
		if err != nil {
			renderError(w, r, fmt.Errorf("can't set media: %w", err))
			return
		}
		renderOK(w, r, map[string]any{})
	}
}

var ErrMediaIsRequired = errors.New("mediafile to set should be provided")
