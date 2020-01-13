package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/alexwilkerson/ddstats-api/pkg/models/postgres"
)

const (
	liveRoom = "live"
)

type PlayerBestReached struct {
	PlayerID         int
	PlayerName       string
	PreviousGameTime float64
}

type PlayerBestSubmitted struct {
	PlayerName       string
	GameID           int
	GameTime         float64
	PreviousGameTime float64
}

type PlayerAbove1000 struct {
	PlayerID   int
	PlayerName string
}

type PlayerDied struct {
	PlayerName string
	GameID     int
	GameTime   float64
	DeathType  string
}

// Hub is the struct which holds the internal communication channels
// for communication with websockets
type Hub struct {
	DB                  *postgres.Postgres
	CurrentID           uint
	Register            chan *Client
	Unregister          chan *Client
	RegisterPlayer      chan *PlayerWithLock
	UnregisterPlayer    chan *PlayerWithLock
	SubmitGame          chan int
	DiscordBroadcast    chan interface{}
	PlayerBestReached   chan *PlayerBestReached
	PlayerBestSubmitted chan *PlayerBestSubmitted
	PlayerAbove1000     chan *PlayerAbove1000
	PlayerDied          chan *PlayerDied
	Players             *sync.Map
	Rooms               map[string]map[*Client]bool
	Broadcast           chan *Message
	quit                chan struct{}
}

// NewHub returns a Hub
func NewHub(db *postgres.Postgres) *Hub {
	rooms := make(map[string]map[*Client]bool)
	rooms[liveRoom] = make(map[*Client]bool)
	return &Hub{
		DB:                  db,
		CurrentID:           1,
		Register:            make(chan *Client),
		Unregister:          make(chan *Client),
		RegisterPlayer:      make(chan *PlayerWithLock),
		UnregisterPlayer:    make(chan *PlayerWithLock),
		SubmitGame:          make(chan int),
		DiscordBroadcast:    make(chan interface{}),
		PlayerBestReached:   make(chan *PlayerBestReached),
		PlayerBestSubmitted: make(chan *PlayerBestSubmitted),
		PlayerAbove1000:     make(chan *PlayerAbove1000),
		PlayerDied:          make(chan *PlayerDied),
		Players:             &sync.Map{},
		Rooms:               rooms,
		Broadcast:           make(chan *Message),
		quit:                make(chan struct{}),
	}
}

// Start is intended to be run in a go routine and will handle all communication
// with websockets.
func (hub *Hub) Start() {
	for {
		select {
		case gameID := <-hub.SubmitGame:
			_ = gameID
			game, err := hub.DB.Games.Get(gameID)
			if err != nil {
				fmt.Println(err)
				break
			}
			for client := range hub.Rooms[liveRoom] {
				message, err := NewMessage(client.Room, "game_submitted", game)
				if err != nil {
					fmt.Println(err)
					break
				}
				err = client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		case player := <-hub.RegisterPlayer:
			hub.Players.Store(player, true)
			for client := range hub.Rooms[liveRoom] {
				message, err := NewMessage(client.Room, "player_logged_in", struct {
					Players []Player `json:"players"`
				}{
					Players: hub.LivePlayers(),
				})
				if err != nil {
					fmt.Println(err)
					break
				}
				err = client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		case player := <-hub.UnregisterPlayer:
			hub.Players.Delete(player)
			for client := range hub.Rooms[liveRoom] {
				message, err := NewMessage(liveRoom, "player_logged_off", struct {
					Players []Player `json:"players"`
				}{
					Players: hub.LivePlayers(),
				})
				if err != nil {
					fmt.Println(err)
					break
				}
				err = client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		case client := <-hub.Register:
			client.ID = hub.CurrentID
			hub.CurrentID++
			if _, ok := hub.Rooms[client.Room]; !ok {
				hub.Rooms[client.Room] = make(map[*Client]bool)
			}
			hub.Rooms[client.Room][client] = true
			fmt.Printf("Size of room %q connections: %d\n", client.Room, len(hub.Rooms[client.Room]))
			if client.Room == liveRoom {
				message, err := NewMessage(liveRoom, "player_list", struct {
					Players []Player `json:"players"`
				}{
					Players: hub.LivePlayers(),
				})
				if err != nil {
					fmt.Println(err)
					break
				}
				err = client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					break
				}
				break
			}
			for client := range hub.Rooms[client.Room] {
				fmt.Println(client.ID)
				message, err := NewMessage(client.Room, "user_connected", struct {
					UserCount int `json:"user_count"`
				}{
					UserCount: len(hub.Rooms[client.Room]),
				})
				if err != nil {
					fmt.Println(err)
					break
				}
				err = client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		case client := <-hub.Unregister:
			delete(hub.Rooms[client.Room], client)
			fmt.Printf("Size of room %q connections: %d\n", client.Room, len(hub.Rooms[client.Room]))
			if len(hub.Rooms[client.Room]) == 0 {
				delete(hub.Rooms, client.Room)
				continue
			}
			for client := range hub.Rooms[client.Room] {
				fmt.Println(client.ID)
				message, err := NewMessage(client.Room, "user_disconnected", struct {
					UserCount int `json:"user_count"`
				}{
					UserCount: len(hub.Rooms[client.Room]),
				})
				if err != nil {
					fmt.Println(err)
					break
				}
				err = client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		case message := <-hub.Broadcast:
			if _, ok := hub.Rooms[message.Room]; !ok {
				break
			}
			fmt.Println("Sending message to all clients in room:", message.Room)
			for client := range hub.Rooms[message.Room] {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		case <-hub.quit:
			return
		}
	}
}

func (hub *Hub) Close() {
	hub.quit <- struct{}{}
}

func toJSONString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), err
}
