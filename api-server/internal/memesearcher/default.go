package memesearcher

import (
	"context"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
	"slices"
	"sort"
	"strings"
)

var _ Engine = &Default{}

type Default struct {
	store models.MemeRepo
}

func newDefault(store models.MemeRepo) *Default {
	return &Default{store: store}
}

// SearchForBoard implements Engine.
func (m *Default) SearchForBoard(ctx context.Context, id models.BoardID, req map[string]string, offset int, limit int, sortBy string) ([]ScoreResult, error) {
	// TODO apply sortBy
	slog := slog.Default().With("from", "API.SearchForBoard")
	slog.InfoContext(ctx, "Started")

	memes := []models.Meme{}
	getOffset := 0
	for {
		newMemes, err := m.store.GetMemesByBoardID(ctx, id, getOffset, 200)
		if err != nil {
			return nil, fmt.Errorf("can't get memes with offset %d: %w", getOffset, err)
		}
		if len(newMemes) == 0 {
			break
		}
		getOffset += len(newMemes)
		memes = append(memes, newMemes...)
	}
	slog.Debug("memes get", "memes", memes)
	res := make([]ScoreResult, 0, len(memes))
	for _, meme := range memes {
		score := m.score(ctx, req, meme.Description)
		if score < 0.1 {
			continue
		}
		res = append(res, ScoreResult{
			Score: score,
			ID:    meme.ID,
		})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Score > res[j].Score
	})

	right := min(offset+limit, len(res))
	return res[offset:right], nil

}

func (m *Default) score(ctx context.Context, req map[string]string, desc map[string]string) float64 {
	s, ok := req["general"]
	if !ok || len(s) == 0 {
		slog.DebugContext(ctx, "Can't find 'genereal' field in request")
		return -1
	}
	words := strings.Split(s, " ")
	matches := map[string]int{}
	for _, ds := range desc {
		for _, word := range strings.Split(ds, " ") {
			if slices.Contains(words, word) {
				matches[word] += 1
			}
		}
	}
	return float64(len(matches)) / float64(len(words))
}
