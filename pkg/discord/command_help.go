package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commandHelp() {
	command := Command{
		name:        "help",
		cooldown:    time.Second,
		description: "List all of my commands or info about a specific command.",
		aliases:     []string{"commands"},
		usage:       "[command name]",
		args:        false,
		getEmbed: func(m *discordgo.MessageCreate, args ...string) *discordgo.MessageEmbed {
			userChannel, err := d.Session.UserChannelCreate(m.Author.ID)
			if err != nil {
				d.errorLog.Printf("error creating user channel for user %s with ID %s: %w", m.Author.Username, m.Author.ID, err)
				return errorEmbed(fmt.Sprintf("Unable to message you. Do you have DMs disabled? %s", m.Author.Mention()))
			}
			if len(args) == 0 {
				commands := d.getCommands()
				_, err = d.Session.ChannelMessageSendEmbed(userChannel.ID, &discordgo.MessageEmbed{
					Description: "You can send `.help [command name]` to get info on a specific command.",
					Color:       defaultColor,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "ddstats.com",
						IconURL: iconURL,
					},
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "List of commands",
							Value: strings.Join(commands, ", "),
						},
					},
				})
				if err != nil {
					d.errorLog.Printf("error sending message to user channel for user %s with ID %s: %w", m.Author.Username, m.Author.ID, err)
					return errorEmbed(fmt.Sprintf("Unable to message you. Do you have DMs disabled? %s", m.Author.Mention()))
				}
				// if the incoming command is coming in through a DM, don't return
				// the notification message that DM was sent to the channel.
				if m.ChannelID == userChannel.ID {
					return nil
				}
				return &discordgo.MessageEmbed{
					Description: fmt.Sprintf("%s I've sent you a DM with all of my commands.", m.Author.Mention()),
					Color:       defaultColor,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "ddstats.com",
						IconURL: iconURL,
					},
				}
			}

			commandName := args[0]
			c, ok := d.commands.Load(commandName)
			if !ok {
				return errorEmbed(fmt.Sprintf("No %q command found. %s", commandName, m.Author.Mention()))
			}
			command := c.(*Command)

			var fields []*discordgo.MessageEmbedField
			if len(command.aliases) > 0 {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:  "Aliases",
					Value: strings.Join(command.aliases, ", "),
				})
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Usage",
				Value: fmt.Sprintf("%s%s %s", prefix, command.name, command.usage),
			})
			cooldown := "None"
			if command.cooldown > 0 {
				cooldown = "1 second"
				if command.cooldown > time.Second {
					cooldown = fmt.Sprintf("%d seconds", int(command.cooldown/time.Second))
				}
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Cooldown",
				Value: cooldown,
			})

			return &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Help %s%s", prefix, command.name),
				Description: command.description,
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
