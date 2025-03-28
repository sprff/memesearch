package statemachine

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Event struct {
	upd tgbotapi.Update
}

type State interface {
	Process(e Event) (State, error)
}
