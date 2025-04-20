package statemachine

import (
	"api-client/pkg/models"
	"fmt"
	"log/slog"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processInline(q *tgbotapi.InlineQuery, r RequestContext) {
	ctx := r.Ctx
	memes, err := r.ApiClient.SearchMemes(ctx, 1, 10, q.Query)
	if err != nil {
		slog.ErrorContext(ctx, "Can't search", "err", err)
		return
	}

	inlineResponse := []any{}
	for _, meme := range memes {
		entry, err := prepareMeme(meme.Meme, r)

		if err != nil {
			slog.ErrorContext(ctx, "Can't prepare meme", "err", err)
			continue
		}
		inlineResponse = append(inlineResponse, entry)
	}
	r.Bot.AnswerInlineQuery(ctx, q.ID, inlineResponse, "")
}

func prepareMeme(meme models.Meme, r RequestContext) (any, error) {
	ctx := r.Ctx
	media, err := r.ApiClient.GetMediaByID(ctx, models.MediaID(meme.ID))
	if err != nil {
		return nil, fmt.Errorf("can't get media: %w", err)
	}
	cm, err := r.Bot.GetFileID(string(meme.ID), path.Ext(meme.Filename), media.Body)
	if err != nil {
		return nil, fmt.Errorf("can't get file id: %w", err)
	}
	switch cm.Type {
	case "photo":
		photo := tgbotapi.NewInlineQueryResultCachedPhoto(string(meme.ID), cm.ID)
		return photo, nil
	case "video":
		video := tgbotapi.NewInlineQueryResultCachedVideo(string(meme.ID), cm.ID, "Title")
		return video, nil
	default:
		return nil, fmt.Errorf("unexpected cm.Type: %s", cm.Type)
	}
}
