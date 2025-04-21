package api

import (
	"context"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
	"memesearch/internal/searchranker"
)

func (a *api) Search(ctx context.Context, req map[string]string, offset, limit int) ([]searchranker.ScroredMeme, error) {
	logger := slog.Default().With("from", "api.SearchMemeByBoardID")
	logger.InfoContext(ctx, "Started")

	batchSize := 200
	memes := []models.Meme{}
	listOffset := 0
	for {
		nmemes, err := a.ListMemes(ctx, listOffset, batchSize, "id")
		if err != nil {
			return nil, fmt.Errorf("can't list memes with offset %d: %w", listOffset, err)
		}
		if len(nmemes) == 0 {
			break
		}
		listOffset += batchSize
		memes = append(memes, nmemes...)
	}
	res, err := a.ranker.Rank(memes, req)
	begin := min(offset, len(res))
	end := min(offset+limit, len(res))

	res = res[begin:end]

	if err != nil {
		return nil, fmt.Errorf("can't rank: %w", err)
	}
	return res, nil
}
