package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	motdMaxLength          = 67
	motdAuthorizedUsername = "vhsx"
)

func (d *Discord) commandMOTD() {
	command := Command{
		name:        "motd",
		cooldown:    5 * time.Second,
		description: "Set the message of the day shown to clients on login. Restricted to the bot owner.",
		usage:       "<message>",
		args:        true,
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			if !strings.EqualFold(m.Author.Username, motdAuthorizedUsername) {
				return errorEmbed(fmt.Sprintf("You are not authorized to use this command. %s", m.Author.Mention()))
			}

			message := strings.TrimSpace(m.Content[len(prefix+"motd"):])
			if message == "" {
				return errorEmbed(fmt.Sprintf("Usage: %smotd %s", prefix, "<message>"))
			}
			if len(message) > motdMaxLength {
				return errorEmbed(fmt.Sprintf("Message is %d characters, max is %d. %s", len(message), motdMaxLength, m.Author.Mention()))
			}

			err := d.DB.MOTD.Insert(message)
			if err != nil {
				d.errorLog.Printf("%v", err)
				return errorEmbed(fmt.Sprintf("Failed to set the message of the day. %s", m.Author.Mention()))
			}

			return &discordgo.MessageEmbed{
				Title:       "Message of the Day Updated",
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
