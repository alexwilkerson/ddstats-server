package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandNews() {
	command := Command{
		name:        "news",
		cooldown:    5 * time.Second,
		description: "Post a new entry to the News feed on the website. Restricted to the bot owner.",
		usage:       "<message>",
		args:        true,
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			if !isBotOwner(m.Author) {
				return errorEmbed(fmt.Sprintf("You are not authorized to use this command. %s", m.Author.Mention()))
			}

			message := strings.TrimSpace(m.Content[len(prefix+"news"):])
			if message == "" {
				return errorEmbed(fmt.Sprintf("Usage: %snews %s", prefix, "<message>"))
			}

			err := d.DB.News.Insert(message)
			if err != nil {
				d.errorLog.Printf("%v", err)
				return errorEmbed(fmt.Sprintf("Failed to post the news entry. %s", m.Author.Mention()))
			}

			return &discordgo.MessageEmbed{
				Title:       "News Posted",
				Description: message,
				Color:       defaultColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
			}
		},
	}
	command.register(d)
}
