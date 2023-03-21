package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

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
		discordBot.SendMessage(channel_id, strings.Replace(message.Content, "\\n", "\n", -1))
	})
	http.ListenAndServe(":3000", nil)
}
