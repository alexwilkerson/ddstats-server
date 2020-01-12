package discord

import (
	"errors"
	"fmt"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/models"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandMe() {
	command := Command{
		name:        "me",
		cooldown:    5 * time.Second,
		description: fmt.Sprintf("Get your personal stats. You must use %sregister [player id] before using the %sme command.", prefix, prefix),
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			d.infoLog.Println(m.Author.ID)
			discordUser, err := d.DB.DiscordUsers.Select(m.Author.ID)
			if err != nil {
				if errors.Is(err, models.ErrNoDiscordUserFound) {
					return &discordgo.MessageEmbed{
						Title:       "Please Register",
						Description: fmt.Sprintf("In order to use the `%sme` command, you must first register using `%sregister [player id]`\n\nTo get your player id, you can use the rank command: `%srank [rank]` or the search command: `%ssearch [player name]`", prefix, prefix, prefix, prefix),
						Color:       defaultColor,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "ddstats.com",
							IconURL: iconURL,
						},
					}
				}
				d.errorLog.Printf("%w", err)
				return errorEmbed(fmt.Sprintf("Database error while trying to retrieve user ID %q. %s", m.Author.ID, m.Author.Mention()))
			}
			player, err := d.ddAPI.UserByID(discordUser.DDID)
			if err != nil {
				if errors.Is(err, ddapi.ErrStatusCode) {
					d.errorLog.Printf("%w", err)
					return errorEmbed(fmt.Sprintf("Unable to access the Devil Daggers API. %s", m.Author.Mention()))
				}
				if errors.Is(err, ddapi.ErrPlayerNotFound) {
					return errorEmbed(fmt.Sprintf("No players were found for Player ID %d. %s", discordUser.DDID, m.Author.Mention()))
				}
				d.errorLog.Printf("%w", err)
				return errorEmbed(fmt.Sprintf("Some error occurred while calling !id. %s", m.Author.Mention()))
			}
			return &discordgo.MessageEmbed{
				Title: fmt.Sprintf("%s (%d)", player.PlayerName, player.PlayerID),
				Color: defaultColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
				Description: fmt.Sprintf("`Rank %d`", player.Rank),
				Fields:      fieldsFromPlayer(player),
			}
		},
	}
	command.register(d)
}
