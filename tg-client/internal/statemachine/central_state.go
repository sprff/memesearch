package statemachine

import (
	"api-client/pkg/models"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"tg-client/internal/telegram"
)

var _ State = &CentralState{}

type CentralState struct {
}

var ErrBadCommandUsage = errors.New("Wrong command usage. Please use help")

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
		case "/login":
			if len(args) < 2 {
				return s, ErrBadCommandUsage
			}
			err := doLogin(r, args[0], args[1])
			return s, err
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
		default:
			r.SendMessage("Unknown command, please use help")
			return s, nil
		}
	case isSearchRequest(r):
		err := doSearchRequest(r)
		return &CentralState{}, err
	case isAddPhoto(r), isAddVideo(r):
		err := doAddMedia(r)
		return &CentralState{}, err
	case isAddVideo(r):
		doAddMedia(r)
		return &CentralState{}, nil
	default:
		return &CentralState{}, nil
	}
}

func isCommand(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	return strings.HasPrefix(r.Event.Message.Text, "/")
}

func isSearchRequest(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) != 0 ||
		msg.Video != nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return msg.Text != ""
}

func doSearchRequest(r RequestContext) error {
	ctx := r.Ctx
	msg := r.Event.Message
	text := msg.Text
	slog.InfoContext(ctx, "doSearchRequest")
	memes, err := r.ApiClient.SearchMemes(ctx, 1, 10, text)
	if err != nil {
		return fmt.Errorf("can't do search request: %w", err)
	}

	if len(memes) == 0 {
		return models.ErrMemeNotFound
	}
	sendMemes(memes, r)
	return nil
}

func sendMemes(memes []models.ScoredMeme, r RequestContext) {
	ctx := r.Ctx
	mges := []telegram.MediaGroupEntry{}
	for _, m := range memes {
		meme := m.Meme
		caption := fmt.Sprintf("ID:%s\nScore:%v\nBoard:%s\nDesc:%s", meme.ID, m.Score, meme.BoardID, meme.Descriptions)
		fileID, err := r.Bot.GetFileID(string(meme.ID), func() ([]byte, error) {
			media, err := r.ApiClient.GetMediaByID(ctx, models.MediaID(meme.ID))
			return media.Body, err
		})
		if err != nil {
			r.SendMessage("Unexpected error")
			slog.ErrorContext(ctx, "Can't get media",
				"error", err.Error(),
				"meme_id", meme.ID)
			continue
		}
		mges = append(mges, telegram.MediaGroupEntry{ID: fileID.ID, Filename: meme.Filename, Caption: caption})
	}
	r.SendMediaGroup(mges)
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
