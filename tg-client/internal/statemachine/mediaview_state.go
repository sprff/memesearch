package statemachine

import (
	"api-client/pkg/models"
	"context"
	"fmt"
	"log/slog"
	"tg-client/internal/telegram"
)

var _ State = &MediaViewState{}

type MediaViewState struct {
	page      int
	getMedias func(ctx context.Context, page int, pageSize int) ([]models.ScoredMeme, error)
	skip      bool
}

// Process implements State.
func (m *MediaViewState) Process(r RequestContext) (State, error) {
	if !m.skip {
		if r.Event.Message == nil {
			return m, nil
		}
		cmd := r.Event.Message.Text
		if cmd == "/next" {
			m.page += 1
		} else if cmd == "/exit" {
			r.SendMessage("Exited media view.")
			return &CentralState{}, nil
		} else {
			r.SendMessage("Use /next to see next results or /exit to exit media view.")
			return m, nil
		}
	}

	m.skip = false
	memes, err := m.getMedias(r.Ctx, m.page, 10)
	if err != nil {
		return &CentralState{}, fmt.Errorf("can't get medias: %w", err)
	}
	sendMemes(r, memes)

	return m, nil
}

func sendMemes(r RequestContext, memes []models.ScoredMeme) {
	if len(memes) == 0 {
		r.SendMessage("No more memes")
	}
	mgs := make([]telegram.MediaGroupEntry, 0, 10)
	for _, m := range memes {
		mge, err := prepareMemeMediaGroup(r, m)
		if err != nil {
			slog.ErrorContext(r.Ctx, "Can't prepare media", "err", err)
			continue
		}
		mgs = append(mgs, mge)
	}
	r.SendMediaGroup(mgs)

}

func prepareMemeMediaGroup(r RequestContext, m models.ScoredMeme) (telegram.MediaGroupEntry, error) {
	ctx := r.Ctx

	meme := m.Meme
	caption := fmt.Sprintf("ID:%s\nScore:%v\nBoard:%s\nDesc:%s", meme.ID, m.Score, meme.BoardID, meme.Descriptions)
	cm, err := r.Bot.GetCachedMedia(string(meme.ID), func() ([]byte, error) {
		media, err := r.ApiClient.GetMediaByID(ctx, models.MediaID(meme.ID))
		return media.Body, err
	})
	if err != nil {
		return telegram.MediaGroupEntry{}, fmt.Errorf("can't get media: %v", err)
	}
	return telegram.MediaGroupEntry{ID: cm.ID, Filename: meme.Filename, Caption: caption, MIMEType: cm.Type}, nil
}
