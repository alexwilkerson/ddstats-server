package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	defaultColor = 0x9A1000 // dark reddish color
)

const (
	numTokensMinimum = 1
)

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages by bot
	if m.Author.Bot {
		return
	}
	contentTokens := strings.Split(strings.TrimSpace(m.Content), " ")
	if len(contentTokens) < numTokensMinimum {
		return
	}

	switch strings.ToLower(contentTokens[0]) {
	case "ping":
		embed := discordgo.MessageEmbed{
			Title: "Amer Ican",
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
