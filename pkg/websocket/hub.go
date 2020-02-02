package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
)

const (
	defaultRoom = "default"
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

type PlayerAboveThreshold struct {
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
	DB               *postgres.Postgres
	CurrentID        uint
	Register         chan *Client
	Unregister       chan *Client
	JoinRoom         chan *Client
	LeaveRoom        chan *Client
	RegisterPlayer   chan *PlayerWithLock
	UnregisterPlayer chan *PlayerWithLock
	SubmitGame       chan int
	DiscordBroadcast chan interface{}
	Players          *sync.Map
	Clients          map[*Client]bool
	Rooms            map[string]map[*Client]bool
	Broadcast        chan *Message
	BroadcastToAll   chan *Message
	quit             chan struct{}
}

// NewHub returns a Hub
func NewHub(db *postgres.Postgres) *Hub {
	rooms := make(map[string]map[*Client]bool)
	rooms[defaultRoom] = make(map[*Client]bool)
	return &Hub{
		DB:               db,
		CurrentID:        1,
		Register:         make(chan *Client, 20),
		Unregister:       make(chan *Client, 20),
		JoinRoom:         make(chan *Client, 20),
		LeaveRoom:        make(chan *Client, 20),
		RegisterPlayer:   make(chan *PlayerWithLock, 20),
		UnregisterPlayer: make(chan *PlayerWithLock, 20),
		SubmitGame:       make(chan int, 20),
		DiscordBroadcast: make(chan interface{}, 20),
		Players:          &sync.Map{},
		Clients:          map[*Client]bool{},
		Rooms:            rooms,
		Broadcast:        make(chan *Message, 20),
		BroadcastToAll:   make(chan *Message, 20),
		quit:             make(chan struct{}),
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
			for client := range hub.Rooms[defaultRoom] {
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
			for client := range hub.Clients {
				message, err := NewMessage(client.Room, "player_logged_in", Player{
					ID:   player.ID,
					Name: player.Name,
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
			for client := range hub.Clients {
				message, err := NewMessage(defaultRoom, "player_logged_off", struct {
					PlayerID int `json:"player_id"`
				}{
					PlayerID: player.ID,
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
			// this ID stuff might be unnecessary
			client.ID = hub.CurrentID
			hub.CurrentID++
			hub.Clients[client] = true
			// if _, ok := hub.Rooms[client.Room]; !ok {
			// 	hub.Rooms[client.Room] = make(map[*Client]bool)
			// }
			// hub.Rooms[client.Room][client] = true
			// fmt.Printf("Size of room %q connections: %d\n", client.Room, len(hub.Rooms[client.Room]))
			// if client.Room == defaultRoom {
			message, err := NewMessage(defaultRoom, "player_list", struct {
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
			// break
			// }
			// if client.Room != defaultRoom {
			// 	for client := range hub.Rooms[client.Room] {
			// 		fmt.Println(client.ID)
			// 		message, err := NewMessage(client.Room, "user_connected", struct {
			// 			UserCount int `json:"user_count"`
			// 		}{
			// 			UserCount: len(hub.Rooms[client.Room]),
			// 		})
			// 		if err != nil {
			// 			fmt.Println(err)
			// 			break
			// 		}
			// 		err = client.Conn.WriteJSON(message)
			// 		if err != nil {
			// 			fmt.Println(err)
			// 			break
			// 		}
			// 	}
			// }
		case client := <-hub.Unregister:
			delete(hub.Clients, client) // delete the client from the list of clients
			// fmt.Printf("Size of room %q connections: %d\n", client.Room, len(hub.Rooms[client.Room]))
			if client.Room != "" { // meaning user does not belong to any room
				if _, ok := hub.Rooms[client.Room]; ok { // verify that the user exists in the room
					delete(hub.Rooms[client.Room], client) // if they do, delete them from the room

					if len(hub.Rooms[client.Room]) == 0 { // if the room is empty, delete the room
						delete(hub.Rooms, client.Room)
						break
					}
					for client := range hub.Rooms[client.Room] { // notify each other client in that room that a player left
						fmt.Println(client.ID)
						message, err := NewMessage(client.Room, "user_count", struct {
							Count int `json:"count"`
						}{
							Count: len(hub.Rooms[client.Room]),
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
				}
			}
		case client := <-hub.JoinRoom: // make sure the room is set when the client calls this function
			if _, ok := hub.Rooms[client.Room]; !ok { // if the room doesn't exist, create it
				hub.Rooms[client.Room] = make(map[*Client]bool)
			}
			hub.Rooms[client.Room][client] = true // add client to room
			fmt.Printf("Size of room %q connections: %d\n", client.Room, len(hub.Rooms[client.Room]))

			for client := range hub.Rooms[client.Room] { // send user count update to each client including this client
				fmt.Println(client.ID)
				message, err := NewMessage(client.Room, "user_count", struct {
					Count int `json:"count"`
				}{
					Count: len(hub.Rooms[client.Room]),
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
		case client := <-hub.LeaveRoom:
			if _, ok := hub.Rooms[client.Room]; ok { // verify that the user exists in the room
				delete(hub.Rooms[client.Room], client) // if they do, delete them from the room
				client.Room = ""                       // make sure their room is empty

				if len(hub.Rooms[client.Room]) == 0 { // if the room is empty, delete the room
					delete(hub.Rooms, client.Room)
					break
				}

				for client := range hub.Rooms[client.Room] { // notify each other client in that room that a player left
					fmt.Println(client.ID)
					message, err := NewMessage(client.Room, "user_count", struct {
						Count int `json:"count"`
					}{
						Count: len(hub.Rooms[client.Room]),
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
			}
		case message := <-hub.Broadcast:
			// a room only exists if a client user is connected to the server
			if _, ok := hub.Rooms[message.Room]; !ok {
				break
			}
			for client := range hub.Rooms[message.Room] {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		case message := <-hub.BroadcastToAll:
			for client := range hub.Clients {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		case <-hub.quit:
			return
		}
	}
}

func (hub *Hub) Close() {
	close(hub.quit)
}

func toJSONString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), err
}
