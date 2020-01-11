package discord

import (
	"testing"
)

func TestStartsWith(t *testing.T) {
	tests := []struct {
		source  string
		pattern string
		want    bool
	}{
		{"!hello", "!", true},
		{"hello", "!", false},
		{"#!bash", "#!", true},
		{"!bash", "#!", false},
		{"!", "!!", false},
		{"!!", "!", true},
		{"", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			got := startsWith(tt.source, tt.pattern)
			if got != tt.want {
				t.Errorf("got %t; want %t", got, tt.want)
			}
		})
	}
}
