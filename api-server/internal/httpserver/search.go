package httpserver

import (
	"fmt"
	"log/slog"
	apiservice "memesearch/internal/api"
	"memesearch/internal/models"
	"net/http"

	"github.com/go-chi/chi"
)

func SearchByBoard(a *apiservice.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := slog.Default().With("from", "Server.SearchByBoard")
		id := models.BoardID(chi.URLParam(r, "id"))
		logger.InfoContext(ctx, "Started", "id", id)

		var req map[string]string
		err := readBody(r, &req)
		if err != nil {
			renderError(w, r, fmt.Errorf("can't read body: %w", err))
			return
		}

		memes, err := a.SearchMemeByBoardID(ctx, id, req, 0, 100)
		if err != nil {
			renderError(w, r, fmt.Errorf("can't search meme by board: %w", err))
			return
		}
		renderOK(w, r, memes)
	}
}
