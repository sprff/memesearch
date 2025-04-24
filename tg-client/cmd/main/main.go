package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"tg-client/internal/contextlogger"
	"tg-client/internal/kvstore"
	"tg-client/internal/statemachine"
	"tg-client/internal/telegram"
	"time"
)

func main() {

	url := os.Getenv("API_SERVER")
	setLogger()
	path := os.Getenv("MS_DATA_FOLDER")
	w, err := NewWrapper(fmt.Sprintf("%s/cache.db", path))
	if err != nil {
		log.Fatalf("Failed to create wrapper: %v", err)
	}

	bot, err := telegram.NewMSBot(os.Getenv("MS_TGCLIENT_BOT_TOKEN"), w)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	state, err := statemachine.New(bot, url, os.Getenv("MS_DATA_FOLDER"))
	if err != nil {
		log.Fatalf("Failed to create statemachine: %v", err)
	}
	state.Process()

}

func setLogger() {
	logFolder := os.Getenv("LOG_FOLDER")
	date := time.Now().Format("2006-01-02")
	handler, err := contextlogger.NewContextHandler(fmt.Sprintf("%s/%s.log", logFolder, date), &slog.HandlerOptions{Level: slog.LevelDebug})
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	slog.SetDefault(slog.New(handler))
}

var _ telegram.CachedMediaStorage = &Wrapper{}

type Wrapper struct {
	s kvstore.Store[telegram.CachedMedia]
}

func NewWrapper(path string) (*Wrapper, error) {
	store, err := kvstore.New[telegram.CachedMedia](path)
	if err != nil {
		return nil, fmt.Errorf("can't create store: %w", err)
	}
	return &Wrapper{
		s: store,
	}, nil
}

// Get implements telegram.CachedMediaStorage.
func (w *Wrapper) Get(ctx context.Context, key string) (telegram.CachedMedia, error) {
	if v, ok := w.s.Get(key); ok {
		return v, nil
	}
	return telegram.CachedMedia{}, fmt.Errorf("enpty")
}

// Set implements telegram.CachedMediaStorage.
func (w *Wrapper) Set(ctx context.Context, key string, value telegram.CachedMedia) error {
	return w.s.Set(key, value)
}
