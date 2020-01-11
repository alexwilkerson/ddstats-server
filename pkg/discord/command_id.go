package discord

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandID() {
	command := Command{
		name:        "id",
		cooldown:    5 * time.Second,
		description: "Retrieve stats for player by their Devil Daggers ID.",
		usage:       "[player id]",
		args:        true,
		aliases:     []string{"score", "pb"},
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return errorEmbed(fmt.Sprintf("Player ID must be an integer. %s", m.Author.Mention()))
			}
			// This clips the args down to a length of 3 because for whatever reason the Devil Daggers
			// API won't find users if the there are more than 2 spaces in their name
			player, err := d.ddAPI.UserByID(id)
			if err != nil {
				if errors.Is(err, ddapi.ErrPlayerNotFound) {
					return errorEmbed(fmt.Sprintf("No players were found for Player ID %d. %s", id, m.Author.Mention()))
				}
				d.errorLog.Printf("%w", err)
				return errorEmbed(fmt.Sprintf("Some error occurred while calling !id. %s", m.Author.Mention()))
			}
			return &discordgo.MessageEmbed{
				Title: fmt.Sprintf("%s (%d)", player.PlayerName, player.PlayerID),
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
				Description: fmt.Sprintf("Rank %d", player.Rank),
				Fields:      fieldsFromPlayer(player),
			}
		},
	}
	command.register(d)
}
