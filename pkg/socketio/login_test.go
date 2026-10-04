package socketio

import (
	"errors"
	"testing"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"
	"github.com/alexwilkerson/ddstats-server/pkg/models"
)

func TestResolveLoginPlayer(t *testing.T) {
	errTimeout := errors.New("context deadline exceeded")
	errDB := errors.New("connection refused")
	ddFound := func(int) (*ddapi.Player, error) {
		return &ddapi.Player{PlayerID: 7, PlayerName: "from dd", GameTime: 1100}, nil
	}
	ddFails := func(err error) func(int) (*ddapi.Player, error) {
		return func(int) (*ddapi.Player, error) { return nil, err }
	}
	dbFound := func(int) (*models.Player, error) {
		return &models.Player{ID: 7, PlayerName: "from db", GameTime: 1050}, nil
	}
	dbFails := func(err error) func(int) (*models.Player, error) {
		return func(int) (*models.Player, error) { return nil, err }
	}

	tests := []struct {
		name       string
		dd         func(int) (*ddapi.Player, error)
		db         func(int) (*models.Player, error)
		wantErr    bool
		wantName   string
		wantBest   float64
		wantFromDD bool
	}{
		{"DD answers", ddFound, dbFails(errDB), false, "from dd", 1100, true},
		{"DD times out, stored player", ddFails(errTimeout), dbFound, false, "from db", 1050, false},
		{"DD says not found, stored player", ddFails(ddapi.ErrPlayerNotFound), dbFound, false, "from db", 1050, false},
		{"DD fails, no stored player", ddFails(errTimeout), dbFails(models.ErrNoRecord), true, "", 0, false},
		{"DD fails, DB fails", ddFails(errTimeout), dbFails(errDB), true, "", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveLoginPlayer(7, tt.dd, tt.db)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("got %+v; want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.id != 7 || got.name != tt.wantName || got.bestGameTime != tt.wantBest {
				t.Errorf("got %+v; want id 7, name %q, best %v", got, tt.wantName, tt.wantBest)
			}
			if (got.ddPlayer != nil) != tt.wantFromDD {
				t.Errorf("got ddPlayer %v; want from DD %t", got.ddPlayer, tt.wantFromDD)
			}
		})
	}
}
