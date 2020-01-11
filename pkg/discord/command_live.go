package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandLive() {
	command := Command{
		name:        "live",
		cooldown:    5 * time.Second,
		description: "Display list of users currently using ddstats.",
		args:        false,
		getEmbed: func(args ...string) *discordgo.MessageEmbed {
			livePlayers := d.websocketHub.LivePlayers()
			var fields []*discordgo.MessageEmbedField
			for _, player := range livePlayers {
				var value string
				if player.Status == "Alive" {
					value = fmt.Sprintf("%s at %.4f seconds.\nLive: https://ddstats.com/user/%d", player.Status, player.GameTime, player.ID)
				} else {
					value = fmt.Sprintf("%s.", player.Status)
				}
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:  fmt.Sprintf("%s (%d)", player.Name, player.ID),
					Value: value,
				})
			}
			var description string
			if len(livePlayers) == 0 {
				description = "No users are live right now."
			}
			return &discordgo.MessageEmbed{
				Title:       "Live Users",
				Description: description,
				Color:       defaultColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
				Fields: fields,
			}
		},
	}
	command.register(d)
}
