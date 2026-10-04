package discord

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	numTokensMinimum = 1
)

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			d.errorLog.Printf("panic handling discord message %q from %s: %v\n%s", m.Content, describeMessageSource(m), r, debug.Stack())
		}
	}()

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
	source := describeMessageSource(m)
	since := time.Since(command.lastUsed)
	if since < command.cooldown {
		d.infoLog.Printf("discord command %s%s from %s rejected: on cooldown", prefix, command.name, source)
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The %s%s command is on cooldown. Please wait %s to use it. %s", prefix, command.name, command.cooldown-since, m.Author.Mention()))
		if err != nil {
			d.errorLog.Printf("sending cooldown notice for %s%s to %s: %v", prefix, command.name, source, err)
		}
		return
	}
	var args []string
	if len(contentTokens) > 1 {
		args = contentTokens[1:]
	}
	command.lastUsed = time.Now()
	embed := command.getEmbed(m, args...)
	d.infoLog.Printf("discord command %s%s args=%q from %s took %s", prefix, command.name, args, source, time.Since(command.lastUsed).Round(time.Millisecond))
	if embed == nil {
		return
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		d.errorLog.Printf("sending %s%s reply to %s: %v", prefix, command.name, source, err)
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
