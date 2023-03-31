package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"log"
)

type WebhookPost struct {
	Content string
}

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	channel_id := os.Getenv("DISCORD_CHANNEL_ID")
	bot_prefix := os.Getenv("DISCORD_BOT_PREFIX")

	if token == "" || channel_id == "" {
		log.Fatal("DISCORD_TOKEN or DISCORD_CHANNEL_ID need to be set.")
		return
	}

	if bot_prefix == "" {
		bot_prefix = "!"
	}

	// setup
	if _, err := os.Stat("tiles"); errors.Is(err, os.ErrNotExist) {
		log.Println("unzipping lfs/tiles.zip to tiles directory")
		err = UnzipFile("lfs/tiles.zip")
		if err != nil {
			log.Println(err)
		}
	}

	var discordBot Bot
	discordBot.Start(token, bot_prefix)

	// webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// print request for reference
		err := PrintWebhook(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var message WebhookPost
		err = json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		eventMessages := ExtractEvents(message.Content)
		for _, message := range eventMessages {
			discordBot.SendMessage(channel_id, fmt.Sprintf("(%s)%s - %s\n", message[1], message[2], message[3]))
		}
	})
	http.ListenAndServe(":3000", nil)
}
