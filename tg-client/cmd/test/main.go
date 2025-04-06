package main

import (
	"api-client/pkg/client"
	"log"
	"os"
	"tg-client/internal/statemachine"
	"tg-client/internal/telegram"
)

func main() {
	bot, err := telegram.NewMSBot(os.Getenv("MS_TGCLIENT_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	client := &client.Client{Url: "http://localhost:1781"}

	state := statemachine.New(bot, client)
	state.Process()

}
