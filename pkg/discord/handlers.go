package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	numTokensMinimum = 1
	prefix           = "."
)

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages by bot
	if m.Author.Bot || !startsWith(m.Content, prefix) {
		return
	}
	contentTokens := strings.Split(strings.TrimSpace(strings.ToLower(m.Content)[len(prefix):]), " ")
	if len(contentTokens) < numTokensMinimum {
		return
	}

	potentialCommand := contentTokens[0]
	c, ok := d.commands.Load(potentialCommand)
	if !ok {
		return
	}
	command := c.(*Command)
	since := time.Since(command.lastUsed)
	if since < command.cooldown {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The %s%s command is on cooldown. Please wait %s to use it. %s", prefix, command.name, command.cooldown-since, m.Author.Mention()))
		if err != nil {
			d.errorLog.Printf("%w", err)
		}
		return
	}
	var args []string
	if len(contentTokens) > 1 {
		args = contentTokens[1:]
	}
	command.lastUsed = time.Now()
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, command.getEmbed(m, args...))
	if err != nil {
		d.errorLog.Printf("%w", err)
	}

	// switch strings.ToLower(contentTokens[0]) {
	// case "ping":
	// 	embed := discordgo.MessageEmbed{
	// 		Title: "Amer Ican",
	// 	}
	// 	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	// case "pong":
	// 	s.ChannelMessageSend(m.ChannelID, "Ping!")
	// }
}
