package httpserver

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	apiservice "memesearch/internal/api"
	"memesearch/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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
			logger.Error("MEDIA IS REQIURED")
			render.JSON(w, r, map[string]any{"status": "MEDIA_IS_REQIURED", "data": map[string]any{}})
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
