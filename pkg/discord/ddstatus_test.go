package discord

import (
	"errors"
	"testing"
	"time"
)

func TestDDStatusRecord(t *testing.T) {
	errDown := errors.New("timeout")
	start := time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)

	type step struct {
		err  error
		want ddStatusTransition
	}
	tests := []struct {
		name      string
		steps     []step
		wantKnown bool
		wantUp    bool
		wantExact bool
		wantSince int
	}{
		{
			name:      "first success establishes up without announcing",
			steps:     []step{{nil, ddStatusNoChange}},
			wantKnown: true, wantUp: true, wantExact: false, wantSince: 0,
		},
		{
			name:      "failures below threshold leave status unknown",
			steps:     []step{{errDown, ddStatusNoChange}, {errDown, ddStatusNoChange}},
			wantKnown: false,
		},
		{
			name:      "failures at threshold at startup establish down without announcing",
			steps:     []step{{errDown, ddStatusNoChange}, {errDown, ddStatusNoChange}, {errDown, ddStatusNoChange}},
			wantKnown: true, wantUp: false, wantExact: false, wantSince: 2,
		},
		{
			name: "up then sustained failures announces down",
			steps: []step{
				{nil, ddStatusNoChange},
				{errDown, ddStatusNoChange},
				{errDown, ddStatusNoChange},
				{errDown, ddStatusWentDown},
				{errDown, ddStatusNoChange},
			},
			wantKnown: true, wantUp: false, wantExact: true, wantSince: 3,
		},
		{
			name: "a success resets the failure count",
			steps: []step{
				{nil, ddStatusNoChange},
				{errDown, ddStatusNoChange},
				{errDown, ddStatusNoChange},
				{nil, ddStatusNoChange},
				{errDown, ddStatusNoChange},
				{errDown, ddStatusNoChange},
			},
			wantKnown: true, wantUp: true, wantExact: false, wantSince: 0,
		},
		{
			name: "down at startup then success announces up",
			steps: []step{
				{errDown, ddStatusNoChange},
				{errDown, ddStatusNoChange},
				{errDown, ddStatusNoChange},
				{nil, ddStatusCameUp},
				{nil, ddStatusNoChange},
			},
			wantKnown: true, wantUp: true, wantExact: true, wantSince: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ddStatus{}
			for i, st := range tt.steps {
				got := s.record(st.err, start.Add(time.Duration(i)*time.Minute))
				if got != st.want {
					t.Fatalf("step %d: got transition %d; want %d", i, got, st.want)
				}
			}
			snap := s.snapshot()
			if snap.known != tt.wantKnown {
				t.Fatalf("got known %t; want %t", snap.known, tt.wantKnown)
			}
			if !tt.wantKnown {
				return
			}
			if snap.up != tt.wantUp {
				t.Errorf("got up %t; want %t", snap.up, tt.wantUp)
			}
			if snap.sinceExact != tt.wantExact {
				t.Errorf("got sinceExact %t; want %t", snap.sinceExact, tt.wantExact)
			}
			wantSince := start.Add(time.Duration(tt.wantSince) * time.Minute)
			if !snap.since.Equal(wantSince) {
				t.Errorf("got since %v; want %v", snap.since, wantSince)
			}
		})
	}
}

func TestDDStatusEmbed(t *testing.T) {
	now := time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		status    ddStatusSnapshot
		wantTitle string
		wantDesc  string
	}{
		{
			name:      "unknown",
			status:    ddStatusSnapshot{},
			wantTitle: "Devil Daggers server status is unknown",
			wantDesc:  "I haven't finished checking the Devil Daggers servers yet. Try again in a minute.",
		},
		{
			name:      "down with exact start",
			status:    ddStatusSnapshot{known: true, up: false, since: now.Add(-90 * time.Minute), sinceExact: true},
			wantTitle: "The Devil Daggers servers are down",
			wantDesc:  "They have been down for 1h 30m.",
		},
		{
			name:      "up since bot startup",
			status:    ddStatusSnapshot{known: true, up: true, since: now.Add(-5 * time.Minute)},
			wantTitle: "The Devil Daggers servers are up",
			wantDesc:  "They have been up for at least 5m.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ddStatusEmbed(tt.status, now)
			if got.Title != tt.wantTitle {
				t.Errorf("got title %q; want %q", got.Title, tt.wantTitle)
			}
			if got.Description != tt.wantDesc {
				t.Errorf("got description %q; want %q", got.Description, tt.wantDesc)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		dur  time.Duration
		want string
	}{
		{30 * time.Second, "less than a minute"},
		{time.Minute, "1m"},
		{59*time.Minute + 59*time.Second, "59m"},
		{2*time.Hour + 5*time.Minute, "2h 5m"},
		{3*24*time.Hour + 4*time.Hour + 7*time.Minute, "3d 4h 7m"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatDuration(tt.dur)
			if got != tt.want {
				t.Errorf("got %q; want %q", got, tt.want)
			}
		})
	}
}
