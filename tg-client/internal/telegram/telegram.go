package telegram

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"tg-client/internal/kvstore"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MSBot struct {
	bot   *tgbotapi.BotAPI
	cache kvstore.Store[CachedMedia]
}

func NewMSBot(token string, datadir string) (*MSBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	cache, err := kvstore.New[CachedMedia](fmt.Sprintf("%s/cache.db", datadir))
	if err != nil {
		return nil, fmt.Errorf("can't create kv: %w", err)
	}
	return &MSBot{
		bot:   bot,
		cache: cache,
	}, nil
}

func (b *MSBot) UpdateChan() tgbotapi.UpdatesChannel {
	return b.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
}

func (b *MSBot) SendMessage(ctx context.Context, chatID int64, text string) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	m, err := b.bot.Send(msg)
	if err != nil {
		slog.ErrorContext(ctx, "can't send message", "error", err.Error(), "chat", chatID, "text", text)
		return 0, fmt.Errorf("can't send message: %w", err)
	}
	return m.MessageID, nil
}

func (b *MSBot) SendMessageReply(ctx context.Context, chatID int64, text string, replyTo int) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = replyTo
	m, err := b.bot.Send(msg)
	if err != nil {
		slog.ErrorContext(ctx, "can't send message", "error", err.Error(), "chat", chatID, "text", text)
		return 0, fmt.Errorf("can't send message: %w", err)
	}
	return m.MessageID, nil
}

type MediaGroupEntry struct {
	ID       string
	Filename string
	Caption  string
	Body     []byte
}

func (b *MSBot) SendMediaGroup(ctx context.Context, chatID int64, medias []MediaGroupEntry) {
	mediaGroup := make([]interface{}, 0, len(medias))
	for _, media := range medias {
		file := tgbotapi.FileID(media.ID)
		switch path.Ext(media.Filename) {
		case ".jpg", ".png":
			inputMedia := tgbotapi.NewInputMediaPhoto(file)
			inputMedia.Caption = media.Caption
			mediaGroup = append(mediaGroup, inputMedia)

		case ".mp4":
			inputMedia := tgbotapi.NewInputMediaVideo(file)
			inputMedia.Caption = media.Caption
			mediaGroup = append(mediaGroup, inputMedia)
		default:
			b.SendMessage(ctx, chatID, "Unexpected file format")
			slog.ErrorContext(ctx, "Unexpected file format",
				"filename", media.Filename)
		}
	}

	msg := tgbotapi.NewMediaGroup(chatID, mediaGroup)
	_, _ = b.bot.Send(msg)

}

func (b *MSBot) GetFileBytes(message *tgbotapi.Message) (string, []byte, error) {
	var fileID string

	if len(message.Photo) > 0 {
		fileID = (message.Photo)[len(message.Photo)-1].FileID
	} else if message.Video != nil {
		fileID = message.Video.FileID
	} else if message.Document != nil {
		fileID = message.Document.FileID
	} else if message.Audio != nil {
		fileID = message.Audio.FileID
	} else if message.Voice != nil {
		fileID = message.Voice.FileID
	} else {
		return "", nil, errors.New("message does not contain a supported file type")
	}
	slog.Info("fileID", "fileID", fileID)
	fileURL, err := b.bot.GetFileDirectURL(fileID)
	if err != nil {
		return "", nil, err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return "", nil, err
	}

	return path.Base(fileURL), buf.Bytes(), nil
}

func (b *MSBot) AnswerInlineQuery(ctx context.Context, inlineQueryID string, results []any, next_offset string) error {
	config := tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		Results:       results,
		NextOffset:    next_offset,
		CacheTime:     1,
	}
	_, err := b.bot.Request(config)
	if err != nil {
		slog.ErrorContext(ctx, "can't answer inline query",
			"error", err.Error(),
			"inlineQueryID", inlineQueryID)
		return fmt.Errorf("can't answer inline query: %w", err)
	}
	return nil
}

type CachedMedia struct {
	ID   string
	Type string // TODO enum

}

func (b *MSBot) GetFileID(key string, getBody func() ([]byte, error)) (res CachedMedia, err error) {
	if cm, ok := b.cache.Get(key); ok {
		return cm, nil
	}
	body, err := getBody()
	if err != nil {
		return CachedMedia{}, fmt.Errorf("can't get body: %w", err)
	}

	sendTo := int64(-1002391398173)
	contentType := http.DetectContentType(body)

	switch contentType {
	case "video/mp4":
		file := tgbotapi.FileBytes{Name: "vieo.mp4", Bytes: body}
		msg := tgbotapi.NewVideo(sendTo, file)
		sentMsg, err := b.bot.Send(msg)
		if err != nil {
			return res, fmt.Errorf("can't send message: %w", err)
		}
		fileID := sentMsg.Video.FileID
		cm := CachedMedia{ID: fileID, Type: "video"}
		b.cache.Set(key, cm)
		return cm, nil
	case "image/png", "image/jpg", "image/jpeg":
		file := tgbotapi.FileBytes{Name: "photo.jpg", Bytes: body}
		msg := tgbotapi.NewPhoto(sendTo, file)
		sentMsg, err := b.bot.Send(msg)
		if err != nil {
			return res, fmt.Errorf("can't send message: %w", err)
		}
		fileID := sentMsg.Photo[len(sentMsg.Photo)-1].FileID
		cm := CachedMedia{ID: fileID, Type: "photo"}
		b.cache.Set(key, cm)
		return cm, nil
	default:
		return res, fmt.Errorf("unexpected file format %s", contentType)
	}

}
