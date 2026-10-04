package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandUp() {
	command := Command{
		name:        "up",
		cooldown:    5 * time.Second,
		description: "Check whether the Devil Daggers servers are up or down, and for how long.",
		args:        false,
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			return ddStatusEmbed(d.ddStatus.snapshot(), time.Now())
		},
	}
	command.register(d)
}
