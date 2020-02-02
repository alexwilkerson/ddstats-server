package websocket

import "sync"

// Player struct is used to represent a player who is currently
// playing the game. The fields are populated from the Devil Daggers
// backend api in the ddapi package
type Player struct {
	ID       int     `json:"player_id"`
	Name     string  `json:"player_name"`
	GameTime float64 `json:"game_time"`
	Status   string  `json:"status"`
}

type PlayerWithLock struct {
	sync.Mutex
	Player
}

func (hub *Hub) LivePlayers() []Player {
	players := []Player{}
	hub.Players.Range(func(k interface{}, v interface{}) bool {
		player := k.(*PlayerWithLock)
		player.Lock()
		players = append(players, player.Player)
		player.Unlock()
		return true
	})
	return players
}
