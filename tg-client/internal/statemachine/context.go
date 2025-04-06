package statemachine

import (
	"api-client/pkg/client"
	"context"
	"tg-client/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RequestContext struct {
	Ctx       context.Context
	Bot       *telegram.MSBot
	Event     *tgbotapi.Update
	ApiClient *client.Client
}

func (r RequestContext) MustChat() int64 {
	chat := r.Event.FromChat()
	if chat == nil {
		panic("expected chat not to be nil")
	}
	return chat.ID
}
