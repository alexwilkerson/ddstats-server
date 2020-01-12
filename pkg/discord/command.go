package discord

import (
	"sort"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	d.commandGlobal()
	d.commandHelp()
	d.commandMe()
	d.commandRegister()
}

func fieldsFromPlayer(player *ddapi.Player) []*discordgo.MessageEmbedField {
	p := message.NewPrinter(language.English)
	return []*discordgo.MessageEmbedField{
		newEmbedField("Time", p.Sprintf("%.4fs", player.GameTime), true),
		newEmbedField("Enemies Killed", p.Sprintf("%d", int(player.EnemiesKilled)), true),
		newEmbedField("Gems", p.Sprintf("%d", int(player.Gems)), true),
		newEmbedField("Accuracy", p.Sprintf("%.2f%%", player.Accuracy), true),
		newEmbedField("Death Type", player.DeathType, true),
		newEmbedField("Overall Time", p.Sprintf("%.4fs", player.OverallTime), true),
		newEmbedField("Overall Time (in days)", p.Sprintf("%.1f days", player.OverallTime/secondsInDay), true),
		newEmbedField("Overall Enemies Killed", p.Sprintf("%d", int(player.OverallEnemiesKilled)), true),
		newEmbedField("Overall Accuracy", p.Sprintf("%.2f%%", player.OverallAccuracy), true),
		newEmbedField("Overall Deaths", p.Sprintf("%d", int(player.OverallDeaths)), true),
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

func (d *Discord) getCommands() []string {
	var commands []string
	d.commands.Range(func(k interface{}, v interface{}) bool {
		if k.(string) == (*v.(*Command)).name {
			commands = append(commands, k.(string))
		}
		return true
	})
	sort.Strings(commands)
	return commands
}
