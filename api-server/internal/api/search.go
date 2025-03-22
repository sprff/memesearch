package api

import (
	"context"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
)

func (a *API) SearchMemeByBoardID(ctx context.Context, id models.BoardID, req map[string]string, offset, limit int) ([]models.Meme, error) {
	logger := slog.Default().With("from", "API.SearchMemeByBoardID")
	logger.InfoContext(ctx, "Started")

	engine := a.searcher.GetMemeEngine("")
	scores, err := engine.SearchForBoard(ctx, id, req, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't search for board: %w", err)
	}
	logger.Debug("Search scores", "scores", scores)
	memes := make([]models.Meme, 0, len(scores))

	for _, score := range scores {
		meme, err := a.storage.MemeRepo.GetMemeByID(ctx, score.ID)
		if err != nil {
			return nil, fmt.Errorf("can't get meme by id %s: %w", score.ID, err)
		}
		memes = append(memes, meme)
	}
	return memes, nil
}
