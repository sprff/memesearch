package statemachine

import (
	"api-client/pkg/client"
	"context"
	"encoding/json"
	"log/slog"
	"tg-client/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Statemachine struct {
	states map[int64]State
	bot    *telegram.MSBot
	client *client.Client
}

func New(bot *telegram.MSBot, client *client.Client) *Statemachine {
	return &Statemachine{
		states: make(map[int64]State),
		bot:    bot,
		client: client,
	}
}

func (s *Statemachine) Process() {
	updates := s.bot.UpdateChan()
	for update := range updates {
		ctx := context.Background() //with update id
		r := RequestContext{
			Ctx:       ctx,
			Bot:       s.bot,
			Event:     &update,
			ApiClient: s.client}
		logUpdate(ctx, update)
		if chat := update.FromChat(); chat != nil {
			if _, ok := s.states[chat.ID]; !ok {
				s.states[chat.ID] = &CentralState{}
			}
			nw, err := s.states[chat.ID].Process(r)
			if err != nil {
				slog.ErrorContext(ctx, "can't process state",
					"err", err.Error(),
					"update", update)
				continue
			}

			s.states[chat.ID] = nw
		}
		if q := update.InlineQuery; q != nil {
			processInline(q, r)
		}
	}
}

func logUpdate(ctx context.Context, u tgbotapi.Update) {
	data, _ := json.Marshal(u)
	slog.InfoContext(ctx, "New update",
		"data", data,
		"from", u.FromChat(),
	)
}
