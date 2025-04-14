package httpserver

import (
	"fmt"
	"log/slog"
	apiservice "memesearch/internal/api"
	"memesearch/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PostMeme(a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := slog.Default().With("from", "Server.PostMeme")
		logger.InfoContext(ctx, "Started")

		var meme models.Meme
		err := readBody(r, &meme)
		if err != nil {
			renderError(w, r, fmt.Errorf("can't read body: %w", err))
			return
		}
		if meme.Description == nil {
			meme.Description = map[string]string{}
		}

		logger.Debug("Body read", "meme", meme)
		id, err := a.CreateMeme(ctx, meme)
		if err != nil {
			renderError(w, r, fmt.Errorf("can't create meme: %w", err))
			return
		}

		renderOK(w, r, map[string]any{"id": id})
	}
}

func GetMemeByID(a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := slog.Default().With("from", "Server.GetMemeByID")
		id := models.MemeID(chi.URLParam(r, "id"))
		logger.InfoContext(ctx, "Started", "id", id)

		meme, err := a.GetMemeByID(ctx, id)
		if err != nil {
			renderError(w, r, fmt.Errorf("can't get meme: %w", err))
			return
		}

		renderOK(w, r, meme)
	}
}

func PutMeme(a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := slog.Default().With("from", "Server.PutMeme")
		id := models.MemeID(chi.URLParam(r, "id"))
		logger.InfoContext(ctx, "Started", "id", id)

		var meme models.Meme
		err := readBody(r, &meme)
		if err != nil {
			renderError(w, r, fmt.Errorf("can't read body: %w", err))
			return
		}

		if meme.Description == nil {
			meme.Description = map[string]string{}
		}
		meme.ID = id

		logger.Debug("Body read", "meme", meme)

		err = a.UpdateMeme(ctx, meme)
		if err != nil {
			renderError(w, r, fmt.Errorf("can't get meme: %w", err))
			return
		}

		renderOK(w, r, meme)
	}
}

func DeleteMeme(a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := slog.Default().With("from", "Server.DeleteMeme")
		id := models.MemeID(chi.URLParam(r, "id"))
		logger.InfoContext(ctx, "Started", "id", id)

		err := a.DeleteMeme(ctx, id)
		if err != nil {
			if err == models.ErrMemeNotFound {
				renderOK(w, r, map[string]any{"ok": false})
				return
			}
			renderError(w, r, fmt.Errorf("can't get meme: %w", err))
		}
		renderOK(w, r, map[string]any{"ok": true})
	}
}
