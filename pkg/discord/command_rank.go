package discord

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandRank() {
	command := Command{
		name:        "rank",
		cooldown:    5 * time.Second,
		description: "Retrieve stats for player by their rank on the Devil Daggers leaderboard.",
		usage:       "[rank]",
		args:        true,
		aliases:     []string{"stats"},
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			if len(args) == 0 {
				return errorEmbed(fmt.Sprintf("No Rank included. %s", m.Author.Mention()))
			}
			rank, err := strconv.Atoi(args[0])
			if err != nil {
				return errorEmbed(fmt.Sprintf("Player Rank must be an integer. %s", m.Author.Mention()))
			}
			// This clips the args down to a length of 3 because for whatever reason the Devil Daggers
			// API won't find users if the there are more than 2 spaces in their name
			player, err := d.ddAPI.UserByRank(rank)
			if err != nil {
				if errors.Is(err, ddapi.ErrStatusCode) {
					d.errorLog.Printf("%w", err)
					return errorEmbed(fmt.Sprintf("Unable to access the Devil Daggers API. %s", m.Author.Mention()))
				}
				if errors.Is(err, ddapi.ErrPlayerNotFound) {
					return errorEmbed(fmt.Sprintf("No players were found for Player Rank %d. %s", rank, m.Author.Mention()))
				}
				d.errorLog.Printf("%w", err)
				return errorEmbed(fmt.Sprintf("Some error occurred while calling !id. %s", m.Author.Mention()))
			}
			return &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("%s (%d)", player.PlayerName, player.PlayerID),
				Description: fmt.Sprintf("`Rank %d`", player.Rank),
				Color:       defaultColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
				Fields: fieldsFromPlayer(player),
			}
		},
	}
	command.register(d)
}
