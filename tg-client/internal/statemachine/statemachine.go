package statemachine

import (
	"api-client/pkg/client"
	"api-client/pkg/models"
	"context"
	"encoding/json"
	"log/slog"
	"tg-client/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Statemachine struct {
	states    map[int64]State     // per CHAT
	userinfo  map[int64]*UserInfo // per USER
	bot       *telegram.MSBot
	clientUrl string
}

type UserInfo struct {
	client      *client.Client
	activeBoard models.BoardID
}

func New(bot *telegram.MSBot, clientUrl string) *Statemachine {
	return &Statemachine{
		states:    make(map[int64]State),
		userinfo:  make(map[int64]*UserInfo),
		bot:       bot,
		clientUrl: clientUrl,
	}
}

func (s *Statemachine) Process() {
	updates := s.bot.UpdateChan()
	for update := range updates {
		ctx := context.Background() //with update id
		s.processUpdate(ctx, update)

	}
}

func (s *Statemachine) processUpdate(ctx context.Context, u tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "Paniced", "reason", r)
		}
	}()
	logUpdate(ctx, u)

	user := u.SentFrom()
	if user == nil {
		return
	}
	if _, ok := s.userinfo[user.ID]; !ok {
		c, err := client.New(s.clientUrl)
		if err != nil {
			slog.ErrorContext(ctx, "Can't create client", "err", err)
			return
		}
		s.userinfo[user.ID] = &UserInfo{client: &c, activeBoard: "default"}
	}

	r := RequestContext{
		Ctx:       ctx,
		Bot:       s.bot,
		Event:     &u,
		ApiClient: s.userinfo[user.ID].client,
		UserInfo:  s.userinfo[user.ID],
	}
	if user := u.SentFrom(); user != nil {

	}

	if chat := u.FromChat(); chat != nil {
		if _, ok := s.states[chat.ID]; !ok {
			s.states[chat.ID] = &CentralState{}
		}
		if _, ok := s.states[chat.ID]; !ok {
			s.states[chat.ID] = &CentralState{}
		}
		nw, err := s.states[chat.ID].Process(r)
		if err != nil {
			slog.ErrorContext(ctx, "can't process state",
				"err", err.Error(),
				"update", u)
			sendError(r, err)
			return
		}

		s.states[chat.ID] = nw
	}
	if q := u.InlineQuery; q != nil {
		processInline(q, r)

	}
}

func logUpdate(ctx context.Context, u tgbotapi.Update) {
	data, _ := json.Marshal(u)
	slog.InfoContext(ctx, "New update",
		"data", data,
		"from", u.FromChat(),
	)
}
