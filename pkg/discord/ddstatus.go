package discord

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	ddStatusCheckInterval = 5 * time.Second
	// A single timeout is common even when the DD servers are healthy, so the
	// servers are only reported down after several consecutive failed checks.
	ddStatusFailureThreshold = 3
)

type ddStatusTransition int

const (
	ddStatusNoChange ddStatusTransition = iota
	ddStatusWentDown
	ddStatusCameUp
)

// ddStatus tracks whether the Devil Daggers servers are up, based on the
// results of periodic heartbeat checks.
type ddStatus struct {
	sync.Mutex
	known      bool
	up         bool
	since      time.Time
	sinceExact bool
	failures   int
}

type ddStatusSnapshot struct {
	known      bool
	up         bool
	since      time.Time
	sinceExact bool
}

// record applies a heartbeat result and reports whether it changed the status.
// The first result after startup only establishes the status: the bot can't
// know when that state began, so it is not announced as a transition.
func (s *ddStatus) record(err error, now time.Time) ddStatusTransition {
	s.Lock()
	defer s.Unlock()

	if err == nil {
		s.failures = 0
		if s.known && s.up {
			return ddStatusNoChange
		}
		wasKnown := s.known
		s.set(true, now, wasKnown)
		if !wasKnown {
			return ddStatusNoChange
		}
		return ddStatusCameUp
	}

	s.failures++
	if s.known && !s.up {
		return ddStatusNoChange
	}
	if s.failures < ddStatusFailureThreshold {
		return ddStatusNoChange
	}
	wasKnown := s.known
	s.set(false, now, wasKnown)
	if !wasKnown {
		return ddStatusNoChange
	}
	return ddStatusWentDown
}

func (s *ddStatus) set(up bool, now time.Time, sinceExact bool) {
	s.known = true
	s.up = up
	s.since = now
	s.sinceExact = sinceExact
}

func (s *ddStatus) snapshot() ddStatusSnapshot {
	s.Lock()
	defer s.Unlock()
	return ddStatusSnapshot{known: s.known, up: s.up, since: s.since, sinceExact: s.sinceExact}
}

func (d *Discord) monitorDDStatus() {
	ticker := time.NewTicker(ddStatusCheckInterval)
	defer ticker.Stop()
	for {
		d.checkDDStatus()
		select {
		case <-ticker.C:
		case <-d.quit:
			return
		}
	}
}

func (d *Discord) checkDDStatus() {
	err := d.ddAPI.Heartbeat()

	var embed *discordgo.MessageEmbed
	// Only transitions are logged; at this check interval, logging every
	// failed heartbeat would flood the log during a long outage.
	switch d.ddStatus.record(err, time.Now()) {
	case ddStatusWentDown:
		d.infoLog.Printf("Devil Daggers servers are down: %v", err)
		embed = &discordgo.MessageEmbed{
			Title:       "The Devil Daggers servers are down",
			Description: "Leaderboards and score submissions are unavailable. I'll post here as soon as they're back up.",
		}
	case ddStatusCameUp:
		d.infoLog.Println("Devil Daggers servers are back up")
		embed = &discordgo.MessageEmbed{
			Title:       "The Devil Daggers servers are back up",
			Description: "Leaderboards and score submissions are available again.",
		}
	default:
		return
	}
	err = d.broadcast(embed)
	if err != nil {
		d.errorLog.Printf("%+v", err)
	}
}

func ddStatusEmbed(status ddStatusSnapshot, now time.Time) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color: defaultColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "ddstats.com",
			IconURL: iconURL,
		},
	}
	if !status.known {
		embed.Title = "Devil Daggers server status is unknown"
		embed.Description = "I haven't finished checking the Devil Daggers servers yet. Try again in a minute."
		return embed
	}

	state := "down"
	if status.up {
		state = "up"
	}
	embed.Title = fmt.Sprintf("The Devil Daggers servers are %s", state)

	duration := formatDuration(now.Sub(status.since))
	if status.sinceExact {
		embed.Description = fmt.Sprintf("They have been %s for %s.", state, duration)
	} else {
		embed.Description = fmt.Sprintf("They have been %s for at least %s.", state, duration)
	}
	return embed
}

// formatDuration renders a duration at minute precision, e.g. "2d 3h 15m".
func formatDuration(dur time.Duration) string {
	if dur < time.Minute {
		return "less than a minute"
	}
	days := int(dur / (24 * time.Hour))
	hours := int(dur % (24 * time.Hour) / time.Hour)
	minutes := int(dur % time.Hour / time.Minute)
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
