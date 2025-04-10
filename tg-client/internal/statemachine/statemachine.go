package statemachine

import (
	"api-client/pkg/client"
	"context"
	"fmt"
	"log/slog"
	"tg-client/internal/telegram"
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
		fmt.Printf("%+v\n", update)
		if chat := update.FromChat(); chat != nil {
			ctx := context.Background()
			if _, ok := s.states[chat.ID]; !ok {
				s.states[chat.ID] = &CentralState{}
			}
			nw, err := s.states[chat.ID].Process(RequestContext{
				Ctx:       ctx,
				Bot:       s.bot,
				Event:     &update,
				ApiClient: s.client})
			if err != nil {
				slog.ErrorContext(ctx, "can't process state",
					"err", err.Error(),
					"update", update)
				continue
			}

			s.states[chat.ID] = nw
		}
		if q := update.InlineQuery; q != nil {
			fmt.Printf("%+v\n", q)
			s.bot.Test(q)
		}
	}
}
