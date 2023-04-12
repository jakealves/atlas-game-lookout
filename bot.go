package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	b.Session.AddHandler(b.reactionHandler)

	err = b.Session.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println("Bot is running !")
}

func (b *Bot) reactionHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	log.Printf("Message reacted: %v, %v. %v", m.UserID, m.Emoji.Name, m.MessageID)
	if string(m.Emoji.Name) == "🧭" {
		// get content from message ID
		message, err := b.Session.ChannelMessage(m.ChannelID, m.MessageID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		label, zoom, long, lat, shift := ParseMapTileCommand(message.Content)
		if long == 0.0 || lat == 0.0 {
			b.SendMessage(m.ChannelID, "I need a long and a lat to generate a tile.")
			return
		}
		generatedTilePath, err := GenerateTileFromCoordinates(label, zoom, long, lat, shift)
		if err != nil {
			log.Printf("Error sending generated tile: %v", err)
		}
		generatedTile, err := os.Open(generatedTilePath)
		if err != nil {
			log.Printf("Error reading generated tile: %v", err)
		}
		defer generatedTile.Close()
		file := &discordgo.File{
			Name:        generatedTilePath,
			ContentType: "image/png",
			Reader:      generatedTile,
		}
		messageEdit := &discordgo.MessageEdit{
			Content: &message.Content,
			Channel: m.ChannelID,
			ID:      m.MessageID,
			Files:   []*discordgo.File{file},
		}
		_, err = s.ChannelMessageEditComplex(messageEdit)
		if err != nil {
			log.Printf("Error sending generated image: %v", err)
		}
	}
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.BotId {
		return
	}
	if m.Content == b.Prefix+"ping" {
		log.Printf("Recieved Ping from %s, sending 'Pong'", m.Author.Username)
		b.SendMessage(m.ChannelID, "pong")
	}
	if strings.Contains(m.Content, b.Prefix+"maptile") {
		log.Printf("Recieved maptile from %s, generating maptile", m.Author.Username)
		label, zoom, long, lat, shift := ParseMapTileCommand(m.Content)

		if long == 0.0 || lat == 0.0 {
			b.SendMessage(m.ChannelID, "I need a long and a lat to generate a tile.")
			return
		}

		generatedTilePath, err := GenerateTileFromCoordinates(label, zoom, long, lat, shift)
		if err != nil {
			log.Printf("Error sending generated tile: %v", err)
		}
		generatedTile, err := os.Open(generatedTilePath)
		if err != nil {
			log.Printf("Error reading generated tile: %v", err)
		}
		defer generatedTile.Close()
		file := &discordgo.File{
			Name:        generatedTilePath,
			ContentType: "image/png",
			Reader:      generatedTile,
		}
		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{file},
		})
		if err != nil {
			log.Printf("Error sending generated image: %v", err)
		}
	}
}

func (b *Bot) SendMessage(channel string, message string) {
	log.Printf("Sending message to %s: %s", channel, message)
	_, _ = b.Session.ChannelMessageSend(channel, message)
}

func ParseMapTileCommand(message string) (label string, zoom int, long float64, lat float64, shift int) {
	var err error

	labelRe := regexp.MustCompile(`(?m)label:"(.+)"`)
	labelMatch := labelRe.FindAllStringSubmatch(message, -1)
	if len(labelMatch) > 0 {
		label = labelMatch[0][1]
	}

	zoomRe := regexp.MustCompile(`(?m)zoom:(\d+)`)
	zoomMatch := zoomRe.FindAllStringSubmatch(message, -1)
	if len(zoomMatch) > 0 {
		zoom, err = strconv.Atoi(zoomMatch[0][1])
		if err != nil {
			log.Printf("Error parsing zoom as int: %v", err)
		}
	}
	if zoom == 0 {
		zoom = 64
	}

	longRe := regexp.MustCompile(`(?im)long:(.+\.\d+) `)
	longSlice := longRe.FindAllStringSubmatch(message, -1)
	if len(longSlice) > 0 {
		long, err = strconv.ParseFloat(strings.TrimSpace(longSlice[0][1]), 64)
		if err != nil {
			log.Printf("Error parsing long as float: %v", err)
		}
	}

	latRe := regexp.MustCompile(`(?im)lat:(.+\.\d+)`)
	latSlice := latRe.FindAllStringSubmatch(message, -1)
	if len(latSlice) > 0 {
		lat, err = strconv.ParseFloat(strings.TrimSpace(latSlice[0][1]), 64)
		if err != nil {
			log.Printf("Error parsing lat as float: %v", err)
		}
	}

	shiftRe := regexp.MustCompile(`(?m)shift:(\d+)`)
	shiftMatch := shiftRe.FindAllStringSubmatch(message, -1)
	if len(shiftMatch) > 0 {
		shift, err = strconv.Atoi(shiftMatch[0][1])
		if err != nil {
			log.Printf("Error parsing shift as int: %v", err)
		}
	}
	return label, zoom, long, lat, shift
}
