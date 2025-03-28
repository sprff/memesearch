package main

import (
	"log"
	"os"
	"tg-client/internal/statemachine"
	"tg-client/internal/telegram"
)

// Пример использования
func main() {
	// Замените "YOUR_BOT_TOKEN" на реальный токен вашего бота
	bot, err := telegram.NewMSBot(os.Getenv("MS_TGCLIENT_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	state := statemachine.New(bot)
	state.Process()

}
