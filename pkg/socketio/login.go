package socketio

import (
	"errors"
	"fmt"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"
	"github.com/alexwilkerson/ddstats-server/pkg/models"
)

// loginPlayer is who a socket.io connection is logged in as.
type loginPlayer struct {
	id           int
	name         string
	bestGameTime float64
	// ddPlayer is set when the DD API answered, and is nil when login fell
	// back to the stored player record.
	ddPlayer *ddapi.Player
	ddErr    error
}

// resolveLoginPlayer looks the player up on the DD API, falling back to the
// stored player record when the DD API fails. dd.hasmodai.com regularly times
// out or returns empty responses, and a failed login closes the client's
// socket; a client that reconnects mid-run then misses its own death
// notification, because game_submitted arrives on a connection that never
// logged in.
func resolveLoginPlayer(id int, ddLookup func(int) (*ddapi.Player, error), dbLookup func(int) (*models.Player, error)) (*loginPlayer, error) {
	dd, ddErr := ddLookup(id)
	if ddErr == nil {
		return &loginPlayer{id: int(dd.PlayerID), name: dd.PlayerName, bestGameTime: dd.GameTime, ddPlayer: dd}, nil
	}

	stored, err := dbLookup(id)
	if errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("DD lookup failed and no stored player: %w", ddErr)
	}
	if err != nil {
		return nil, fmt.Errorf("DD lookup failed (%v), loading stored player: %w", ddErr, err)
	}
	return &loginPlayer{id: stored.ID, name: stored.PlayerName, bestGameTime: stored.GameTime, ddErr: ddErr}, nil
}
