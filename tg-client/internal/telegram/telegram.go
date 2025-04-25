package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"tg-client/internal/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MSBot struct {
	bot        *tgbotapi.BotAPI
	cache      CachedMediaStorage
	uploadChan chan uploadChanEntry
}

func NewMSBot(token string, cache CachedMediaStorage) (*MSBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	mbot := &MSBot{
		bot:        bot,
		cache:      cache,
		uploadChan: make(chan uploadChanEntry, uploadChanSize),
	}
	go mbot.startUploading()
	return mbot, nil
}

func (b *MSBot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	return b.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
}

func (b *MSBot) SendMessage(ctx context.Context, chatID int64, text string, parse *string, replyTo *int, replyKeyboard any) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if parse != nil {
		msg.ParseMode = *parse
	}
	if replyTo != nil {
		msg.ReplyToMessageID = *replyTo
	}
	if replyKeyboard != nil {
		msg.ReplyMarkup = replyKeyboard
	}
	m, err := b.bot.Send(msg)
	if err != nil {
		return 0, fmt.Errorf("can't send message: %w", err)
	}
	return m.MessageID, nil
}

func (b *MSBot) SendMediaGroup(ctx context.Context, chatID int64, entries []MediaGroupEntry) (ids []int, err error) {
	if len(entries) > 10 {
		panic("expected <= 10 medias per group")
	}

	mediaGroup := make([]any, 0, len(entries))
	for _, entry := range entries {
		file := tgbotapi.FileID(entry.Media.FileID)
		switch entry.Media.Type {
		case CMPhoto:
			inputMedia := tgbotapi.NewInputMediaPhoto(file)
			inputMedia.Caption = entry.Caption
			mediaGroup = append(mediaGroup, inputMedia)
		case CMVideo:
			inputMedia := tgbotapi.NewInputMediaVideo(file)
			inputMedia.Caption = entry.Caption
			mediaGroup = append(mediaGroup, inputMedia)
		default:
			panic(fmt.Sprintf("uncovered CMType: %v", entry.Media.Type))
		}
	}

	msg := tgbotapi.NewMediaGroup(chatID, mediaGroup)
	resp, err := b.bot.Request(msg)
	if err != nil {
		return nil, fmt.Errorf("can't request: %w", err)
	}
	var messages []tgbotapi.Message
	err = json.Unmarshal(resp.Result, &messages)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal: %w", err)
	}

	for _, m := range messages {
		ids = append(ids, m.MessageID)
	}
	return ids, nil
}

func (b *MSBot) GetFile(ctx context.Context, message *tgbotapi.Message) (name string, body []byte, err error) {
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
	fileURL, err := b.bot.GetFileDirectURL(fileID)
	if err != nil {
		return "", nil, fmt.Errorf("can't get file url: %w", err)
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", nil, fmt.Errorf("can't get: %w", err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return "", nil, fmt.Errorf("can't copy: %w", err)
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

type uploadChanEntry struct {
	key string
	f   func() (UploadEntry, error)
}

func (b *MSBot) Upload(ctx context.Context, key string, forceUpload bool, getUpload func() (UploadEntry, error)) (res CachedMedia, err error) {
	if cm, err := b.cache.Get(ctx, key); err == nil && !forceUpload {
		return cm, nil
	}
	ph := CachedMedia{FileID: uploadPlaceholderPhoto, Type: CMPhoto}

	//TODO it will be better to check if key is already in queue
	ue := uploadChanEntry{
		key: key,
		f:   getUpload,
	}

	select {
	case b.uploadChan <- ue:
	case <-time.Tick(20 * time.Millisecond):
		//if queue is full just return ph
		return ph, nil
	}

	return ph, nil
}
func (b *MSBot) DeleteMessage(ctx context.Context, chatID int64, msgID int) (err error) {
	msg := tgbotapi.NewDeleteMessage(chatID, msgID)
	_, err = b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("can't delete message: %w", err)
	}
	return nil

}

const (
	uploadChat             int64 = -1002391398173
	uploadChanSize               = 400
	uploadRateLimit              = 1
	uploadRateInterval           = 1 * time.Second
	uploadPlaceholderPhoto       = "AgACAgIAAyEGAASOidcdAAOHaAlkN4uc9yN_v2ikFNWmbF-_JfYAAp74MRtYfEhIu5koMOS_BhsBAAMCAAN5AAM2BA" //TODO config it
)

func (b *MSBot) startUploading() {
	for ue := range b.uploadChan {
		upload, err := ue.f()
		if err != nil {
			slog.Error("Can't get upload body", "err", err)
		}
		if upload.Body != nil {
			err := tgRetry(3, func() error {
				return uploadBody(b, ue.key, upload.Name, *upload.Body)
			})
			if err != nil {
				slog.Error("Can't upload body", "err", err)
			}
		}
	}
}

func tgRetry(retryCnt int, f func() error) error {
	id := utils.GenereateUUIDv7()
	for {
		err := f()
		if err == nil {
			return nil
		}
		var tgerr *tgbotapi.Error
		if errors.As(err, &tgerr) {
			if strings.Contains(tgerr.Message, "Too Many Requests") {
				dur := min(5, tgerr.RetryAfter)
				slog.Info("tg retry sleep", "duration", dur)
				time.Sleep(time.Duration(dur) * time.Second)
				continue
			}
		}
		slog.Warn("Can't do", "retriesRemains", retryCnt, "err", err, "id", id)
		retryCnt -= 1
		if retryCnt == 0 {
			return fmt.Errorf("retries exceded: %s", id)
		}

	}
}

func uploadBody(b *MSBot, key string, name string, body []byte) error {
	contentType := http.DetectContentType(body)

	cm := CachedMedia{}
	switch contentType {
	case "video/mp4":
		file := tgbotapi.FileBytes{Name: fmt.Sprintf("%s.mp4", name), Bytes: body}
		msg := tgbotapi.NewVideo(uploadChat, file)

		sentMsg, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("can't send to uploadChat: %w", err)
		}

		cm.FileID = sentMsg.Video.FileID
		cm.Type = CMVideo
	case "image/png", "image/jpg", "image/jpeg":
		file := tgbotapi.FileBytes{Name: fmt.Sprintf("%s.jpg", name), Bytes: body}
		msg := tgbotapi.NewPhoto(uploadChat, file)
		sentMsg, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("can't send to uploadChat: %w", err)
		}

		cm.FileID = sentMsg.Photo[len(sentMsg.Photo)-1].FileID
		cm.Type = CMPhoto

	default:
		return fmt.Errorf("unexpected file format %s", contentType)
	}

	err := b.cache.Set(context.Background(), key, cm)
	if err != nil {
		return fmt.Errorf("can't set cache: %w", err)
	}
	return nil
}
