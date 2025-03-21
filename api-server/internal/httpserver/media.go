package httpserver

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	apiservice "memesearch/internal/api"
	"memesearch/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func GetMedia(ctx context.Context, a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.Default().With("from", "Server.PutMedia")
		id := models.MediaID(chi.URLParam(r, "id"))
		logger.InfoContext(ctx, "Started", "id", id)

		media, err := a.GetMedia(ctx, models.MediaID(id))
		if err != nil {
			logger.Error("CANT_GET_media") // TODO fix errors in media funcs
			render.JSON(w, r, map[string]any{"status": "CANT_GET_media", "data": map[string]any{}})
			return
		}

		if _, err := io.Copy(w, bytes.NewBuffer(media.Body)); err != nil {
			logger.Error("CANT_Copy_media") // TODO fix errors in media funcs
			render.JSON(w, r, map[string]any{"status": "CANT_Copy_media", "data": map[string]any{}})
			return
		}
	}
}

func PutMedia(ctx context.Context, a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			logger.Error("Can't read file", "error", err.Error())
			render.JSON(w, r, map[string]any{"status": "UNEXPECTED_ERROR", "data": map[string]any{}})
			return
		}
		err = a.SetMedia(ctx, models.Media{ID: id, Body: body})
		if err != nil {
			logger.Error("Can't set media", "error", err.Error())
			render.JSON(w, r, map[string]any{"status": "UNEXPECTED_ERROR", "data": map[string]any{}})
			return
		}
		render.JSON(w, r, map[string]any{"status": "OK", "data": map[string]any{}})
	}
}
