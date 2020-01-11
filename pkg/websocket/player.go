package websocket

import "sync"

// Player struct is used to represent a player who is currently
// playing the game. The fields are populated from the Devil Daggers
// backend api in the ddapi package
type Player struct {
	sync.Mutex
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	GameTime float64 `json:"game_time"`
	Status   string  `json:"status"`
}

func (hub *Hub) LivePlayers() []Player {
	return toPlayerSlice(hub.Players)
}

func toPlayerSlice(m *sync.Map) []Player {
	players := []Player{}
	m.Range(func(k interface{}, v interface{}) bool {
		players = append(players, *k.(*Player))
		return true
	})
	return players
}
