package discord

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/models"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandRegister() {
	command := Command{
		name:        "register",
		cooldown:    5 * time.Second,
		description: fmt.Sprintf("Registers a user using their Devil Daggers ID, in order to use the %sme command.", prefix),
		usage:       "[player id]",
		args:        true,
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			if len(args) == 0 {
				return errorEmbed(fmt.Sprintf("Player ID required. %s", m.Author.Mention()))
			}
			ddID, err := strconv.Atoi(args[0])
			if err != nil {
				d.errorLog.Printf("%+v", err)
				return errorEmbed(fmt.Sprintf("Player ID must be an integer. %s", m.Author.Mention()))
			}
			err = d.DB.DiscordUsers.Upsert(m.Author.ID, ddID)
			if err != nil {
				if errors.Is(err, models.ErrDiscordUserVerified) {
					return errorEmbed(fmt.Sprintf("You are verified. You cannot change your registered Player ID. %s", m.Author.Mention()))
				}
				d.errorLog.Printf("%+v", err)
				return errorEmbed(fmt.Sprintf("Error registering user to database. %s", m.Author.Mention()))
			}
			return &discordgo.MessageEmbed{
				Title:       "User Registered",
				Description: fmt.Sprintf("User `%s` is now registered with DD User ID `%d`. %s", m.Author.Username, ddID, m.Author.Mention()),
				Color:       defaultColor,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "ddstats.com",
					IconURL: iconURL,
				},
			}
		},
	}
	command.register(d)
}
