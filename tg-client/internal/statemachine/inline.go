package statemachine

import (
	"api-client/pkg/models"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"tg-client/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processInline(q *tgbotapi.InlineQuery, r RequestContext) {
	ctx := r.Ctx

	page := 1
	if q.Offset != "" {
		var err error
		page, err = strconv.Atoi(q.Offset)
		if err != nil {
			// очень жаль
		}
	}
	req := q.Query
	filter := telegram.CMPhoto
	if len(req) > 0 && req[0] == '!' {
		filter = telegram.CMVideo
		req = req[1:]

	}

	memes, err := r.ApiClient.SearchMemes(ctx, (page-1)*50, 50, req)
	if err != nil {
		slog.ErrorContext(ctx, "Can't search", "err", err)
		return
	}
	newOffset := ""
	if len(memes) != 0 {
		newOffset = fmt.Sprintf("%d", page+1)
	}
	inlineResponse := []any{}
	for _, meme := range memes {
		entry, err := prepareMeme(meme.Meme, r, filter)

		if err != nil {
			if err != ErrSkipped {
				slog.ErrorContext(ctx, "Can't prepare meme", "err", err)
			}
			continue
		}
		inlineResponse = append(inlineResponse, entry)
	}

	r.Bot.AnswerInlineQuery(ctx, q.ID, inlineResponse, newOffset)
}

var ErrSkipped = errors.New("skip")

func prepareMeme(meme models.Meme, r RequestContext, filter telegram.CachedMediaType) (any, error) {
	ctx := r.Ctx

	cm, err := r.Bot.Upload(ctx, string(meme.ID), false, func() (telegram.UploadEntry, error) {
		media, err := r.ApiClient.GetMediaByID(ctx, models.MediaID(meme.ID))
		if err != nil {
			return telegram.UploadEntry{}, fmt.Errorf("can't get media: %w", err)
		}
		return telegram.UploadEntry{Name: "file", Body: &media.Body}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't get file id: %w", err)
	}
	if cm.Type != filter {
		return nil, ErrSkipped
	}
	switch cm.Type {
	case telegram.CMPhoto:
		photo := tgbotapi.NewInlineQueryResultCachedPhoto(string(meme.ID), cm.FileID)
		return photo, nil
	case telegram.CMVideo:
		video := tgbotapi.NewInlineQueryResultCachedVideo(string(meme.ID), cm.FileID, " ")
		video.Description = meme.Descriptions["general"]
		return video, nil
	default:
		return nil, fmt.Errorf("unexpected cm.Type: %s", cm.Type)
	}
}
