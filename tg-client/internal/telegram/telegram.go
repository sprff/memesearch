package telegram

import (
	"bytes"
	"errors"
	"io"
	"net/http"

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

func (b *MSBot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.bot.Send(msg)
	return err
}

func (b *MSBot) SendPhoto(chatID int64, photo []byte, caption string) error {
	file := tgbotapi.FileBytes{Name: "photo.jpg", Bytes: photo}
	msg := tgbotapi.NewPhoto(chatID, file)
	msg.Caption = caption
	_, err := b.bot.Send(msg)
	return err
}

func (b *MSBot) SendPhotoGroup(chatID int64, photos [][]byte, caption string) error {
	if len(photos) == 0 {
		return errors.New("no photos provided")
	}

	mediaGroup := make([]interface{}, 0, len(photos))
	for i, photo := range photos {
		file := tgbotapi.FileBytes{Name: "photo.jpg", Bytes: photo}
		inputMedia := tgbotapi.NewInputMediaPhoto(file)
		if i == 0 {
			inputMedia.Caption = caption
		}
		mediaGroup = append(mediaGroup, inputMedia)
	}

	msg := tgbotapi.NewMediaGroup(chatID, mediaGroup)
	_, err := b.bot.Send(msg)
	return err
}

func (b *MSBot) SendVideo(chatID int64, video []byte, caption string) error {
	file := tgbotapi.FileBytes{Name: "video.mp4", Bytes: video}
	msg := tgbotapi.NewVideo(chatID, file)
	msg.Caption = caption
	_, err := b.bot.Send(msg)
	return err
}

func (b *MSBot) SendVideoGroup(chatID int64, videos [][]byte, caption string) error {
	if len(videos) == 0 {
		return errors.New("no videos provided")
	}

	mediaGroup := make([]interface{}, 0, len(videos))
	for i, video := range videos {
		file := tgbotapi.FileBytes{Name: "video.mp4", Bytes: video}
		inputMedia := tgbotapi.NewInputMediaVideo(file)
		if i == 0 {
			inputMedia.Caption = caption
		}
		mediaGroup = append(mediaGroup, inputMedia)
	}

	msg := tgbotapi.NewMediaGroup(chatID, mediaGroup)
	_, err := b.bot.Send(msg)
	return err
}

func (b *MSBot) GetFileBytes(message *tgbotapi.Message) ([]byte, error) {
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
		return nil, errors.New("message does not contain a supported file type")
	}

	fileURL, err := b.bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b.SendMessage(message.From.ID, "Start Copy")
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
