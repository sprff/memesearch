package main

import (
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

	state := statemachine.New(bot, "http://localhost:1781")
	state.Process()
}
