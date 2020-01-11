package discord

import (
	"errors"
	"fmt"
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
			userName := strings.Join(args, " ")
			if len(userName) < 2 {
				return errorEmbed(fmt.Sprintf("Username must be longer than 2 characters. %s", m.Author.Mention()))
			}
			players, err := d.ddAPI.UserSearch(userName)
			if err != nil {
				if errors.Is(err, ddapi.ErrNoPlayersFound) {
					return errorEmbed(fmt.Sprintf("No players were found for '%s'. %s", strings.Join(args, " "), m.Author.Mention()))
				}
				d.errorLog.Printf("%w", err)
				return errorEmbed(fmt.Sprintf("Some error occurred while calling !search. %s", m.Author.Mention()))
			}
			player := players[0]
			if len(players) == 1 {
				return &discordgo.MessageEmbed{
					Title:       fmt.Sprintf("%s (%d)", player.PlayerName, player.PlayerID),
					Description: fmt.Sprintf("Rank %d", player.Rank),
					Color:       defaultColor,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "ddstats.com",
						IconURL: iconURL,
					},
					Fields: fieldsFromPlayer(player),
				}
			}
			numPlayersFound := len(players)
			var fieldName string
			if numPlayersFound == 100 {
				fieldName = "Found 100+ users. (showing 10)"
			} else if numPlayersFound > 10 {
				fieldName = fmt.Sprintf("Found %d users. (showing 10)", numPlayersFound)
			} else {
				fieldName = fmt.Sprintf("Found %d users.", numPlayersFound)
			}
			usersString := "```"
			for _, player := range players {
				usersString += fmt.Sprintf("\nName: %s\nRank: %d\nID:   %d\n", player.PlayerName, player.Rank, player.PlayerID)
			}
			usersString += "```"
			return &discordgo.MessageEmbed{
				Title: fmt.Sprintf("User Search: %s", strings.Join(args, " ")),
				Color: defaultColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
				Fields: []*discordgo.MessageEmbedField{
					newEmbedField(fieldName, usersString, false),
				},
			}
		},
	}
	command.register(d)
}
