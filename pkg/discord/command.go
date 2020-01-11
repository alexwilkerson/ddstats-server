package discord

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
	"github.com/bwmarrin/discordgo"
)

const (
	defaultColor = 0x9A1000 // dark reddish color
	iconURL      = "https://ddstats.com/static/img/logo_red_100x100.png"
)

const (
	secondsInDay = 86400
)

type Command struct {
	name        string
	description string
	usage       string
	cooldown    time.Duration
	lastUsed    time.Time
	aliases     []string
	args        bool
	getEmbed    func(*discordgo.MessageCreate, ...string) *discordgo.MessageEmbed
}

func (c *Command) register(d *Discord) {
	d.commands.Store(c.name, c)
	for _, name := range c.aliases {
		d.commands.Store(name, c)
	}
}

func (d *Discord) registerCommands() {
	d.commandTop()
	d.commandLive()
	d.commandSearch()
	d.commandID()
	d.commandRank()
}

func fieldsFromPlayer(player *ddapi.Player) []*discordgo.MessageEmbedField {
	return []*discordgo.MessageEmbedField{
		newEmbedField("Time", fmt.Sprintf("%.4fs", player.GameTime), true),
		newEmbedField("Enemies Killed", strconv.Itoa(int(player.EnemiesKilled)), true),
		newEmbedField("Gems", strconv.Itoa(int(player.Gems)), true),
		newEmbedField("Accuracy", fmt.Sprintf("%.2f%%", player.Accuracy), true),
		newEmbedField("Death Type", player.DeathType, true),
		newEmbedField("Overall Time", fmt.Sprintf("%.4fs", player.OverallTime), true),
		newEmbedField("Overall Time (in days)", fmt.Sprintf("%.1f days", player.OverallTime/secondsInDay), true),
		newEmbedField("Overall Enemies Killed", strconv.Itoa(int(player.OverallEnemiesKilled)), true),
		newEmbedField("Overall Accuracy", fmt.Sprintf("%.2f%%", player.OverallAccuracy), true),
		newEmbedField("Overall Deaths", strconv.Itoa(int(player.OverallDeaths)), true),
	}
}

func newEmbedField(name, value string, inline bool) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
}

func errorEmbed(e string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: e,
		Color:       defaultColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "ddstats.com",
			IconURL: iconURL,
		},
	}
}
