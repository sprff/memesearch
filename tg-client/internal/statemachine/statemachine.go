package statemachine

import (
	"api-client/pkg/client"
	"context"
	"fmt"
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
		if chat := update.FromChat(); chat != nil {
			if _, ok := s.states[chat.ID]; !ok {
				s.states[chat.ID] = &DefaultState{}
			}
			nw, err := s.states[chat.ID].Process(RequestContext{
				Ctx:       context.Background(),
				Bot:       s.bot,
				Event:     &update,
				ApiClient: s.client})
			if err != nil {
				//Log
				s.bot.SendMessage(chat.ID, fmt.Sprintf("Error: %v", err))
				continue
			}
			s.states[chat.ID] = nw
		}
	}
}
