package discord

import (
	"errors"
	"fmt"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (d *Discord) commandGlobal() {
	command := Command{
		name:        "global",
		cooldown:    5 * time.Second,
		description: "Retrieve the global stats from the Devil Daggers database.",
		args:        false,
		aliases:     []string{"server"},
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			leaderboard, err := d.ddAPI.GetLeaderboard(0, 0)
			if err != nil {
				if errors.Is(err, ddapi.ErrStatusCode) {
					d.errorLog.Printf("%w", err)
					return errorEmbed(fmt.Sprintf("Unable to access the Devil Daggers API. %s", m.Author.Mention()))
				}
				d.errorLog.Printf("%w", err)
				return errorEmbed(fmt.Sprintf("Some error occurred while calling !id. %s", m.Author.Mention()))
			}
			p := message.NewPrinter(language.English)
			return &discordgo.MessageEmbed{
				Title: "Global Stats",
				Color: defaultColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Total Players",
						Value:  p.Sprintf("%d", leaderboard.GlobalPlayerCount),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Game Time",
						Value:  p.Sprintf("%.4fs", leaderboard.GlobalGameTime),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Average Game Time",
						Value:  p.Sprintf("%.4fs", leaderboard.GlobalAverageGameTime),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Gems Collected",
						Value:  p.Sprintf("%d", leaderboard.GlobalGems),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Enemies Killed",
						Value:  p.Sprintf("%d", leaderboard.GlobalEnemiesKilled),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Deaths",
						Value:  p.Sprintf("%d", leaderboard.GlobalDeaths),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Daggers Hit",
						Value:  p.Sprintf("%d", leaderboard.GlobalDaggersHit),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Daggers Fired",
						Value:  p.Sprintf("%d", leaderboard.GlobalDaggersFired),
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Global Accuracy",
						Value:  p.Sprintf("%.2f%%", leaderboard.GlobalAccuracy),
						Inline: true,
					},
				},
			}
		},
	}
	command.register(d)
}
