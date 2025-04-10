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

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MSBot struct {
	bot *tgbotapi.BotAPI
}

func NewMSBot(token string) (*MSBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &MSBot{bot: bot}, nil
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

func (b *MSBot) SendError(ctx context.Context, chatID int64, msg string) {
	_, _ = b.SendMessage(ctx, chatID, fmt.Sprintf("error: %s\nrequest-id: <code>%s</code>", msg, "id"))
}

type MediaGroupEntry struct {
	ID       string
	Filename string
	Caption  string
	Body     []byte
}

func (b *MSBot) SendMediaGroup(ctx context.Context, chatID int64, medias []MediaGroupEntry) error {
	mediaGroup := make([]interface{}, 0, len(medias))
	for _, media := range medias {
		file := tgbotapi.FileBytes{Name: media.Filename, Bytes: media.Body}
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
			b.SendError(ctx, chatID, "Unexpected file format")
			slog.ErrorContext(ctx, "Unexpected file format",
				"filename", media.Filename)
		}
	}

	msg := tgbotapi.NewMediaGroup(chatID, mediaGroup)
	_, err := b.bot.Send(msg)
	return err
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

func (b *MSBot) Test(q *tgbotapi.InlineQuery) {
	res := []any{}

	for i := range 2 {
		fileURL := "https://storage.yandexcloud.net/meme-search-local/064db12a-af5e-4a02-b35d-7ede062a1122?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=YCAJEnRQ0arrTvsih5QDhQ_P_%2F20250410%2Fru-central1%2Fs3%2Faws4_request&X-Amz-Date=20250410T181156Z&X-Amz-Expires=600&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=8fedc7732dce242a77baeaa4a33810e556f39dccacbdaf53d811bee98d3d2ae4"
		slog.Info("fileURL", "fileURL", fileURL)
		photo := tgbotapi.NewInlineQueryResultPhoto(fmt.Sprintf("test_%d", i), fileURL)
		photo.ThumbURL = fileURL
		photo.Title = "Title"
		photo.Description = "LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG DESC"
		res = append(res, photo)
		// video := tgbotapi.NewInlineQueryResultVideo(fmt.Sprintf("test_video_%d", i), fileURL)
		// video.ThumbURL = "https://avatars.mds.yandex.net/i?id=10ccd7d2b5e15699ec1ee14eb62a60fb_l-5220614-images-thumbs&n=13"
		// video.MimeType = "video/mp4"
		// video.Title = "Title"
		// res = append(res, video)

	}

	// video := tgbotapi.NewInlineQueryResultVideo("3", fileURL)
	// video.ThumbURL = "https://avatars.mds.yandex.net/i?id=10ccd7d2b5e15699ec1ee14eb62a60fb_l-5220614-images-thumbs&n=13"
	// video.MimeType = "video/mp4"
	// video.Description = "LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG DESC"
	// video.Title = "Title"
	b.AnswerInlineQuery(context.Background(), q.ID, res, "50")

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
