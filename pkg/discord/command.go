package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	defaultColor = 0x9A1000 // dark reddish color
	iconURL      = "https://ddstats.com/static/img/logo_red_100x100.png"
)

type Command struct {
	name        string
	description string
	usage       string
	cooldown    time.Duration
	lastUsed    time.Time
	aliases     []string
	args        bool
	getEmbed    func(*discordgo.MessageCreate, ...string) *discordgo.MessageEmbed
}

func (c *Command) register(d *Discord) {
	d.commands.Store(c.name, c)
	for _, name := range c.aliases {
		d.commands.Store(name, c)
	}
}

func (d *Discord) registerCommands() {
	d.commandTop()
	d.commandLive()
	d.commandSearch()
}

func newEmbedField(name, value string, inline bool) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
}

func errorEmbed(e string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: e,
		Color:       defaultColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "ddstats.com",
			IconURL: iconURL,
		},
	}
}
