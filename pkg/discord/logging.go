package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// discordgoLogger routes discordgo's internal logs into the server's loggers.
// Without this, gateway failures such as a reconnect loop that never succeeds
// are invisible, because discordgo otherwise logs nothing below LogError.
func discordgoLogger(infoLog, errorLog *log.Logger) func(msgL, caller int, format string, a ...interface{}) {
	return func(msgL, caller int, format string, a ...interface{}) {
		logger := infoLog
		if msgL <= discordgo.LogWarning {
			logger = errorLog
		}
		logger.Printf("discordgo: %s", fmt.Sprintf(format, a...))
	}
}

func (d *Discord) onReady(s *discordgo.Session, r *discordgo.Ready) {
	d.infoLog.Printf("Discord bot connected as %s to %d guild(s)", r.User.Username, len(r.Guilds))
}

func (d *Discord) onResumed(s *discordgo.Session, r *discordgo.Resumed) {
	d.infoLog.Println("Discord bot resumed gateway session")
}

func (d *Discord) onDisconnect(s *discordgo.Session, e *discordgo.Disconnect) {
	d.errorLog.Println("Discord bot disconnected from gateway; discordgo will try to reconnect")
}

func describeMessageSource(m *discordgo.MessageCreate) string {
	if m.GuildID == "" {
		return fmt.Sprintf("%s (%s) in DM", m.Author.Username, m.Author.ID)
	}
	return fmt.Sprintf("%s (%s) in guild %s channel %s", m.Author.Username, m.Author.ID, m.GuildID, m.ChannelID)
}
