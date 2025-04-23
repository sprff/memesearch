package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"tg-client/internal/contextlogger"
	"tg-client/internal/statemachine"
	"tg-client/internal/telegram"
	"time"
)

func main() {
	logFolder := os.Getenv("LOG_FOLDER")
	date := time.Now().Format("2006-01-02")
	handler, err := contextlogger.NewContextHandler(fmt.Sprintf("%s/%s.log", logFolder, date), &slog.HandlerOptions{Level: slog.LevelDebug})
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	slog.SetDefault(slog.New(handler))

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
