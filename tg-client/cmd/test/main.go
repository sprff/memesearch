package main

import (
	"api-client/pkg/client"
	"fmt"
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
	client, err := client.New("http://localhost:1781")
	if err != nil {
		fmt.Println("can't get client", err)
	}

	state := statemachine.New(bot, client)
	state.Process()

}
