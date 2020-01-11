package discord

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandSearch() {
	command := Command{
		name:        "search",
		cooldown:    5 * time.Second,
		description: "Search the Devil Daggers database using a search string to find a user by their name on the leaderboard.",
		usage:       "[search string]",
		args:        true,
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			var userName string
			// This clips the args down to a length of 3 because for whatever reason, the Devil Daggers
			// API won't find users if the length is over 3
			if len(args) > 3 {
				userName = strings.Join(args[:3], " ")
			} else {
				userName = strings.Join(args, " ")
			}
			if len(userName) < 2 {
				return errorEmbed(fmt.Sprintf("Username must be longer than 2 characters. %s", m.Author.Mention()))
			}
			players, err := d.ddAPI.UserSearch(userName)
			if err != nil {
				if errors.Is(err, ddapi.ErrNoPlayersFound) {
					return errorEmbed(fmt.Sprintf("No players were found for '%s'. This maybe be because the Devil Daggers API doesn't handle spaces ' ' well. %s", userName, m.Author.Mention()))
				}
				d.errorLog.Printf("%w", err)
				return errorEmbed(fmt.Sprintf("Some error occurred while calling !search. %s", m.Author.Mention()))
			}
			if len(players) == 1 {
				player := players[0]
				return &discordgo.MessageEmbed{
					Title:       fmt.Sprintf("%s (%d)", player.PlayerName, player.PlayerID),
					Description: fmt.Sprintf("Rank %d", player.Rank),
					Color:       defaultColor,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "ddstats.com",
						IconURL: iconURL,
					},
					Fields: []*discordgo.MessageEmbedField{
						newEmbedField("Time", fmt.Sprintf("%.4fs", player.GameTime), true),
						newEmbedField("Enemies Killed", strconv.Itoa(int(player.EnemiesKilled)), true),
					},
				}
			}
			return &discordgo.MessageEmbed{
				Title: fmt.Sprintf("User Search: %s", userName),
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
				Fields: []*discordgo.MessageEmbedField{
					newEmbedField(fmt.Sprintf("Found %d users.", len(players)), userName, false),
				},
			}
		},
	}
	command.register(d)
}
