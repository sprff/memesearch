package memesearcher

import (
	"context"
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"memesearch/internal/storage/psql"
)

type ScoreResult struct {
	Score float64
	ID    models.MemeID
}

type Engine interface {
	SearchForBoard(ctx context.Context, id models.BoardID, req map[string]string, offset, limit int) ([]ScoreResult, error)
}

type Searcher struct {
	memeEngines       map[string]Engine
	defaultMemeEngine string
}

func (s Searcher) GetMemeEngine(memeEngine string) Engine {
	if ms, ok := s.memeEngines[memeEngine]; ok {
		return ms
	}
	return s.memeEngines[s.defaultMemeEngine]
}

func New(cfg config.Config) (Searcher, error) {
	memeStroe, err := psql.NewMemeStore(cfg.Database)
	if err != nil {
		return Searcher{}, fmt.Errorf("can't create memestore: %w", err)
	}
	memeEngines := map[string]Engine{
		"default": newDefault(memeStroe),
	}
	return Searcher{memeEngines: memeEngines, defaultMemeEngine: "default"}, nil
}
