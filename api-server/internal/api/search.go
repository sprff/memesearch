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

	isEmpty := true
	if len(req) > 0 {
		for _, v := range req {
			if len(v) != 0 {
				isEmpty = false
				break
			}
		}
	}

	if isEmpty {
		memes, err := a.ListMemes(ctx, offset, limit, "id")
		if err != nil {
			return nil, fmt.Errorf("can't list memes: %w", err)
		}
		smemes := make([]searchranker.ScroredMeme, 0, len(memes))
		for _, m := range memes {
			smemes = append(smemes, searchranker.ScroredMeme{Score: 0, Meme: m})
		}
		return smemes, nil
	}

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
