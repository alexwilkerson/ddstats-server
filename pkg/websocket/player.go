package websocket

import "sync"

// Player struct is used to represent a player who is currently
// playing the game. The fields are populated from the Devil Daggers
// backend api in the ddapi package
type Player struct {
	ID   int
	Name string
}

func toPlayerSlice(m *sync.Map) []*Player {
	var players []*Player
	m.Range(func(k interface{}, v interface{}) bool {
		player, ok := m.Load(k)
		if !ok {
			return false
		}
		players = append(players, player.(*Player))
		return true
	})
	return players
}
