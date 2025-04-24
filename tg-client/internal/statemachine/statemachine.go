package statemachine

import (
	"api-client/pkg/client"
	"api-client/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"tg-client/internal/contextlogger"
	"tg-client/internal/kvstore"
	"tg-client/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Statemachine struct {
	states   map[int64]State         // per CHAT
	userinfo kvstore.Store[UserInfo] // per USER
	bot      *telegram.MSBot
	client   client.Client
}

type UserInfo struct {
	ActiveBoard models.BoardID
	Token       string
}

func New(bot *telegram.MSBot, clientUrl string, datadir string) (*Statemachine, error) {
	client, err := client.New(clientUrl)
	if err != nil {
		return nil, fmt.Errorf("can't create client: %w", err)
	}
	userinfo, err := kvstore.New[UserInfo](fmt.Sprintf("%s/userinfo.db", datadir))
	if err != nil {
		return nil, fmt.Errorf("can't create userinfo: %w", err)
	}
	return &Statemachine{
		states:   make(map[int64]State),
		userinfo: userinfo,
		bot:      bot,
		client:   client,
	}, nil
}

func (s *Statemachine) Process() {
	updates := s.bot.GetUpdatesChan()
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
	userID := strconv.FormatInt(user.ID, 10)
	if _, ok := s.userinfo.Get(userID); !ok {
		s.userinfo.Set(userID, UserInfo{ActiveBoard: "default"})
	}
	info, _ := s.userinfo.Get(userID)
	c := s.client.WithToken(info.Token)
	rid := c.GenerateID()

	r := RequestContext{
		Ctx:       contextlogger.AppendCtx(ctx, slog.String("request_id", rid)),
		Bot:       s.bot,
		Event:     &u,
		ApiClient: &c,
		UserInfo:  &info,
	}
	defer func(i *UserInfo) {
		s.userinfo.Set(userID, *i)
	}(&info)

	if chat := u.FromChat(); chat != nil {
		if _, ok := s.states[chat.ID]; !ok {
			s.states[chat.ID] = &CentralState{}
		}
		if _, ok := s.states[chat.ID]; !ok {
			s.states[chat.ID] = &CentralState{}
		}
		nw, err := s.states[chat.ID].Process(r)
		if err != nil {
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
