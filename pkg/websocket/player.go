package websocket

import "sync"

// Player struct is used to represent a player who is currently
// playing the game. The fields are populated from the Devil Daggers
// backend api in the ddapi package
type Player struct {
	ID   int
	Name string
}

func (hub *Hub) LivePlayers() []*Player {
	return toPlayerSlice(hub.Players)
}

func toPlayerSlice(m *sync.Map) []*Player {
	players := []*Player{}
	m.Range(func(k interface{}, v interface{}) bool {
		players = append(players, k.(*Player))
		return true
	})
	return players
}
