package main

import (
	"api-client/pkg/client"
	"log"
	"log/slog"
	"os"
	"tg-client/internal/statemachine"
	"tg-client/internal/telegram"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	bot, err := telegram.NewMSBot(os.Getenv("MS_TGCLIENT_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	client := &client.Client{Url: "http://localhost:1781"}

	state := statemachine.New(bot, client)
	state.Process()

}
