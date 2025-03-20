package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	apiservice "memesearch/internal/api"
	"memesearch/internal/models"
	"net/http"

	"github.com/go-chi/chi"
)

func PostMeme() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request, ctx context.Context, a *apiservice.API) (any, error) {
		logger := slog.Default().With("from", "Server.PostMeme")
		logger.InfoContext(ctx, "Started")

		var meme models.Meme
		err := readBody(r, &meme)
		if err != nil {
			return nil, fmt.Errorf("can't read body: %w", err)
		}

		if meme.Descriptions == nil {
			meme.Descriptions = map[string]string{}
		}

		logger.Debug("Read meme", "meme", meme)

		id, err := a.CreateMeme(ctx, meme)
		if err != nil {
			return nil, fmt.Errorf("can't get meme: %w", err)
		}

		return map[string]any{"id": id}, nil
	}
}

func GetMemeByID() handlerWithError {
	return func(w http.ResponseWriter, r *http.Request, ctx context.Context, a *apiservice.API) (any, error) {
		logger := slog.Default().With("from", "Server.GetMemeByID")
		id := chi.URLParam(r, "id")
		logger.InfoContext(ctx, "Started", "id", id)

		meme, err := a.GetMemeByID(ctx, models.MemeID(id))
		if err != nil {
			return nil, fmt.Errorf("can't get meme: %w", err)
		}

		return meme, nil
	}
}
