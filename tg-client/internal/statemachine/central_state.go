package statemachine

import (
	"api-client/pkg/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var _ State = &CentralState{}

type CentralState struct {
}

var ErrBadCommandUsage = errors.New("Wrong command usage. Please use /help")

func (s *CentralState) Process(r RequestContext) (State, error) {
	switch {
	case isCommand(r):
		spl := strings.Split(r.Event.Message.Text, " ")
		cmd := spl[0]
		args := spl[1:]
		switch cmd {
		case "/register":
			if len(args) < 2 {
				return s, ErrBadCommandUsage
			}
			err := doRegister(r, args[0], args[1])
			return s, err

		case "/help":
			r.SendMessage(help())
			return s, nil
		case "/login":
			if len(args) < 2 {
				return s, ErrBadCommandUsage
			}
			err := doLogin(r, args[0], args[1])
			return s, err
		case "/logout":
			r.UserInfo.Token = ""
			r.UserInfo.ActiveBoard = "default"
			r.SendMessage("Logged out")
			return s, nil
		case "/whoami":
			err := doWhoami(r)
			return s, err
		case "/setboard":
			if len(args) < 1 {
				return s, ErrBadCommandUsage
			}
			err := doSetBoard(r, models.BoardID(args[0]))
			return s, err
		case "/getboard":
			err := doGetBoard(r)
			return s, err
		case "/createboard":
			if len(args) < 1 {
				return s, ErrBadCommandUsage
			}
			err := doCreateBoard(r, args[0])
			return s, err
		case "/listboards":
			err := doListBoards(r)
			return s, err
		case "/subscribe":
			if len(args) < 1 {
				return s, ErrBadCommandUsage
			}
			err := doSubscribe(r, models.BoardID(args[0]))
			return s, err
		case "/unsubscribe":
			if len(args) < 1 {
				return s, ErrBadCommandUsage
			}
			err := doUnsubscribe(r, models.BoardID(args[0]))
			return s, err
		case "/search":
			text := strings.Join(args[:], " ")
			mv := MediaViewState{page: 1, skip: true, getMedias: func(ctx context.Context, page, pageSize int) ([]models.ScoredMeme, error) {
				return r.ApiClient.SearchMemes(ctx, page, pageSize, text)
			}}
			return mv.Process(r)
		default:
			r.SendMessage("Unknown command, please use /help")
			return s, nil
		}
	case isAddPhoto(r), isAddVideo(r):
		err := doAddMedia(r)
		return &CentralState{}, err
	default:
		r.SendMessage("Please use /help to see what i can do.")
		return &CentralState{}, nil
	}
}

func isCommand(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	return strings.HasPrefix(r.Event.Message.Text, "/")
}

func isAddPhoto(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) == 0 ||
		msg.Video != nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return true
}

func isAddVideo(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) != 0 ||
		msg.Video == nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return true
}

func doAddMedia(r RequestContext) error {
	ctx := r.Ctx
	msg := r.Event.Message
	filename, media, err := r.Bot.GetFileBytes(msg)
	if err != nil {
		return fmt.Errorf("can't get files: %w", err)

	}

	meme, err := r.ApiClient.PostMeme(ctx, r.UserInfo.ActiveBoard, filename, map[string]string{"general": msg.Caption})
	if err != nil {
		return fmt.Errorf("can't create meme: %w", err)
	}
	err = r.ApiClient.PutMediaByID(ctx, models.Media{ID: models.MediaID(meme.ID), Body: media}, filename)
	if err != nil {
		return fmt.Errorf("can't set media: %w", err)
	}
	slog.InfoContext(ctx, "Meme created",
		"id", meme.ID)
	r.SendMessageReply(fmt.Sprintf("<code>%s</code>", meme.ID), msg.MessageID)
	return nil
}
