package statemachine

import (
	"fmt"
	"tg-client/internal/telegram"
)

type Statemachine struct {
	states map[int64]State
	bot    *telegram.MSBot
}

func New(bot *telegram.MSBot) *Statemachine {
	return &Statemachine{
		states: make(map[int64]State),
		bot:    bot,
	}
}

func (s *Statemachine) Process() {
	updates := s.bot.UpdateChan()
	for update := range updates {
		if chat := update.FromChat(); chat != nil {
			if _, ok := s.states[chat.ID]; !ok {
				s.states[chat.ID] = &DefaultState{bot: s.bot}
			}
			nw, err := s.states[chat.ID].Process(Event{upd: update})
			if err != nil {
				//Log
				s.bot.SendMessage(chat.ID, fmt.Sprintf("Error: %v", err))
				continue
			}
			s.states[chat.ID] = nw
		}
	}
}
