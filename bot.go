package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	BotId   string
	Prefix  string
	Session *discordgo.Session
}

func (b *Bot) Start(token string, prefix string) {
	var err error
	b.Prefix = prefix
	b.Session, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := b.Session.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	b.BotId = u.ID
	b.Session.AddHandler(b.messageHandler)

	err = b.Session.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println("Bot is running !")
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.BotId {
		return
	}
	if m.Content == b.Prefix+"ping" {
		log.Printf("Recieved Ping from %s, sending 'Pong'", m.Author.Username)
		b.SendMessage(m.ChannelID, "pong")
	}
}

func (b *Bot) SendMessage(channel string, message string) {
	log.Printf("Sending message to %s: %s", channel, message)
	_, _ = b.Session.ChannelMessageSend(channel, message)
}
