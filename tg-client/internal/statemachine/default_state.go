package statemachine

import "tg-client/internal/telegram"

var _ State = &DefaultState{}

type DefaultState struct {
	bot *telegram.MSBot
}

func (d *DefaultState) Process(e Event) (State, error) {
	chat := *(e.upd.FromChat())
	d.bot.SendMessage(chat.ID, "Default")
	return d, nil
}
