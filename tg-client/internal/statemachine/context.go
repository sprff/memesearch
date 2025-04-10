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

func (r RequestContext) SendMessage(text string) (int, error) {
	return r.Bot.SendMessage(r.Ctx, r.MustChat(), text)
}

func (r RequestContext) SendMessageReply(text string, replyTo int) (int, error) {
	return r.Bot.SendMessageReply(r.Ctx, r.MustChat(), text, replyTo)
}

func (r RequestContext) SendError(msg string) {
	r.Bot.SendError(r.Ctx, r.MustChat(), msg)
}

func (r RequestContext) SendMediaGroup(medias []telegram.MediaGroupEntry) error {
	return r.Bot.SendMediaGroup(r.Ctx, r.MustChat(), medias)
}
