package discord

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestDiscordgoLogger(t *testing.T) {
	tests := []struct {
		name      string
		level     int
		wantError bool
	}{
		{"error goes to error log", discordgo.LogError, true},
		{"warning goes to error log", discordgo.LogWarning, true},
		{"informational goes to info log", discordgo.LogInformational, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var info, errs bytes.Buffer
			logf := discordgoLogger(log.New(&info, "", 0), log.New(&errs, "", 0))
			logf(tt.level, 0, "reconnect failed: %s", "boom")

			got, other := info.String(), errs.String()
			if tt.wantError {
				got, other = other, got
			}
			if !strings.Contains(got, "discordgo: reconnect failed: boom") {
				t.Errorf("expected message in target log, got %q", got)
			}
			if other != "" {
				t.Errorf("expected nothing in other log, got %q", other)
			}
		})
	}
}

func TestDescribeMessageSource(t *testing.T) {
	author := &discordgo.User{ID: "42", Username: "vhsx"}
	tests := []struct {
		name string
		msg  *discordgo.Message
		want string
	}{
		{"direct message", &discordgo.Message{Author: author, ChannelID: "c1"}, "vhsx (42) in DM"},
		{"guild channel", &discordgo.Message{Author: author, ChannelID: "c1", GuildID: "g1"}, "vhsx (42) in guild g1 channel c1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := describeMessageSource(&discordgo.MessageCreate{Message: tt.msg})
			if got != tt.want {
				t.Errorf("got %q; want %q", got, tt.want)
			}
		})
	}
}
