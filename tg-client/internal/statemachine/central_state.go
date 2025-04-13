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
	case isAddPhoto(r):
		doAddMedia(r)
		return &CentralState{}, nil
	case isAddVideo(r):
		doAddMedia(r)
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
	memes, err := r.ApiClient.SearchMemeByBoardID(ctx, "board", 1, text)
	if err != nil {
		slog.ErrorContext(ctx, "can't do serch request",
			"error", err.Error())
		r.SendMessage("can't do search request")
	}
	sendMemes(memes, r)
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
				r.SendError("Meme don't have media")
				slog.WarnContext(ctx, "Meme don't have media",
					"meme_id", meme.ID)
			default:
				r.SendError("Unexpected error")
				slog.ErrorContext(ctx, "Can't get media",
					"error", err.Error(),
					"meme_id", meme.ID)
			}
			continue
		}
		mges = append(mges, telegram.MediaGroupEntry{Filename: meme.Filename, Caption: caption, Body: media.Body})
	}
	r.SendMediaGroup(mges)
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

func doAddMedia(r RequestContext) {
	ctx := r.Ctx
	msg := r.Event.Message
	filename, media, err := r.Bot.GetFileBytes(msg)
	if err != nil {
		slog.ErrorContext(ctx, "can't get file bytes",
			"err", err.Error(),
			"msg", msg)
		return
	}

	id, err := r.ApiClient.PostMeme(ctx, models.Meme{
		BoardID:      "board",
		Filename:     filename,
		Descriptions: map[string]string{"general": msg.Caption},
	})
	if err != nil {
		r.SendError("can't create meme")
		slog.ErrorContext(ctx, "can't create meme",
			"error", err.Error())
		return
	}
	err = r.ApiClient.PutMedia(ctx, models.Media{ID: models.MediaID(id), Body: media}, filename)
	if err != nil {
		r.SendError("can't set media")
		slog.ErrorContext(ctx, "can't set media",
			"error", err.Error())
		return
	}
	slog.InfoContext(ctx, "Meme created",
		"id", id)
	r.SendMessageReply(fmt.Sprintf("<code>%s</code>", id), msg.MessageID)
}
