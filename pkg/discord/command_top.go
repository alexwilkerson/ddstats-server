package discord

import (
	"errors"
	"fmt"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"
	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandTop() {
	command := Command{
		name:        "top",
		cooldown:    5 * time.Second,
		description: "Show the top 10 players of the Devil Daggers leaderboard.",
		aliases:     []string{"top10", "top-ten"},
		args:        false,
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			leaderboard, err := d.ddAPI.GetLeaderboard(10, 0)
			if err != nil {
				if errors.Is(err, ddapi.ErrStatusCode) {
					d.errorLog.Printf("%w", err)
					return errorEmbed(fmt.Sprintf("Unable to access the Devil Daggers API. %s", m.Author.Mention()))
				}
				d.errorLog.Printf("%w", err)
				return nil
			}
			// fields := make([]*discordgo.MessageEmbedField, 10)
			var fields []*discordgo.MessageEmbedField
			for _, player := range leaderboard.Players {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:  player.PlayerName,
					Value: fmt.Sprintf("%.4f", player.GameTime),
				})
			}
			return &discordgo.MessageEmbed{
				Title: "Top 10",
				Color: defaultColor,
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
