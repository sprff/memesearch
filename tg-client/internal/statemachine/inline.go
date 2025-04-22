package statemachine

import (
	"api-client/pkg/models"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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
	spl := strings.Split(q.Query, "!")
	req := spl[0]
	flags := "pv"
	if len(spl) > 1 {
		req = spl[1]
		flags = spl[0]
	}

	memes, err := r.ApiClient.SearchMemes(ctx, page, 10, req)
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
		entry, err := prepareMeme(meme.Meme, r, flags)

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

func prepareMeme(meme models.Meme, r RequestContext, flags string) (any, error) {
	ctx := r.Ctx
	photos := strings.Contains(flags, "p")
	videos := strings.Contains(flags, "v")

	cm, err := r.Bot.GetCachedMedia(string(meme.ID), func() ([]byte, error) {
		media, err := r.ApiClient.GetMediaByID(ctx, models.MediaID(meme.ID))
		return media.Body, err
	})
	if err != nil {
		return nil, fmt.Errorf("can't get file id: %w", err)
	}
	switch cm.Type {
	case "photo":
		if !photos {
			return nil, ErrSkipped
		}
		photo := tgbotapi.NewInlineQueryResultCachedPhoto(string(meme.ID), cm.ID)
		return photo, nil
	case "video":
		if !videos {
			return nil, ErrSkipped
		}
		video := tgbotapi.NewInlineQueryResultCachedVideo(string(meme.ID), cm.ID, "Title")
		return video, nil
	default:
		return nil, fmt.Errorf("unexpected cm.Type: %s", cm.Type)
	}
}
