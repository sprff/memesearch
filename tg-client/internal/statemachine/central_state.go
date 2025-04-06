package statemachine

import (
	"api-client/pkg/models"
	"errors"
	"fmt"
	"log/slog"
	"tg-client/internal/telegram"
)

var _ State = &CentralState{}

type CentralState struct {
}

func (d *CentralState) Process(r RequestContext) (State, error) {
	switch {
	case isSearchRequest(r):
		doSearchRequest(r)
		return &CentralState{}, nil
	case isInlineSearchRequest(r):
		doInlineSearchRequest(r)
		return &CentralState{}, nil
	case isAddPhoto(r):
		doAddPhoto(r)
		return &CentralState{}, nil
	case isAddVideo(r):
		doAddVideo(r)
		return &CentralState{}, nil
	default:
		return &CentralState{}, nil
	}
}

func isSearchRequest(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) != 0 ||
		msg.Video != nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return msg.Text != ""
}

func doSearchRequest(r RequestContext) {
	ctx := r.Ctx
	msg := r.Event.Message
	text := msg.Text
	slog.InfoContext(ctx, "doSearchRequest")
	memes, err := r.ApiClient.SearchMemeByBoardID(ctx, "board", map[string]string{"general": text})
	if err != nil {
		slog.ErrorContext(ctx, "can't do serch request", "error", err.Error())
		r.Bot.SendMessage(ctx, r.MustChat(), "can't do search request")
	}
	sendMemes(memes, r)
	// for _, meme := range memes {
	// 	sendMeme(meme, r)
	// }

}

func sendMemes(memes []models.Meme, r RequestContext) {
	ctx := r.Ctx
	mges := []telegram.MediaGroupEntry{}
	for _, meme := range memes {
		caption := fmt.Sprintf("ID:%s\nBoard:%s\nDesc:%s", meme.ID, meme.BoardID, meme.Descriptions)
		media, err := r.ApiClient.GetMedia(ctx, models.MediaID(meme.ID))
		if err != nil {
			switch {
			case errors.Is(err, models.ErrMediaNotFound):
				r.Bot.SendError(ctx, r.MustChat(), "Meme don't have media")
				slog.WarnContext(ctx, "meme don't have media", "id", meme.ID)
			default:
				r.Bot.SendError(ctx, r.MustChat(), "Unexpected error")
				slog.ErrorContext(ctx, "can't get media", "error", err.Error(), "id", meme.ID)
			}
			continue
		}
		mges = append(mges, telegram.MediaGroupEntry{Filename: meme.Filename, Caption: caption, Body: media.Body})
	}
	r.Bot.SendMediaGroup(ctx, r.MustChat(), mges)
}

func isInlineSearchRequest(r RequestContext) bool {
	return false
}

func doInlineSearchRequest(r RequestContext) {
	ctx := r.Ctx
	slog.InfoContext(ctx, "doInlineSearchRequest")
}

func isAddPhoto(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) == 0 ||
		msg.Video != nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return true
}

func doAddPhoto(r RequestContext) {
	ctx := r.Ctx
	slog.InfoContext(ctx, "doAddPhoto")
	r.Bot.SendMessage(ctx, r.MustChat(), "doAddPhoto")
}

func isAddVideo(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) != 0 ||
		msg.Video == nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return true
}

func doAddVideo(r RequestContext) {
	ctx := r.Ctx
	slog.InfoContext(ctx, "doAddVideo")
	r.Bot.SendMessage(ctx, r.MustChat(), "doAddVideo")
}
