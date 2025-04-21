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
	bot, err := telegram.NewMSBot(os.Getenv("MS_TGCLIENT_BOT_TOKEN"), os.Getenv("MS_DATA_FOLDER"))
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	state, err := statemachine.New(bot, "http://localhost:1781", os.Getenv("MS_DATA_FOLDER"))
	if err != nil {
		log.Fatalf("Failed to create statemachine: %v", err)
	}
	state.Process()
}
