package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client represents the user connected through the websocket
type Client struct {
	ID   uint
	Conn *websocket.Conn
	Hub  *Hub
	Room string
}

// Message represents the message that will be sent to the client
// and frontend. The message stores the intended "room" that the
// message is intended for, the Type, which is always 1 meaning text,
// a Func name, which is a function name sent to the frontend to be handled,
// and a Body which holds any extra data the function might need to send
type Message struct {
	Room string `json:"-"`
	Type int    `json:"type,omitempty"`
	Func string `json:"func,omitempty"`
	Body string `json:"body,omitempty"`
}

// NewMessage populates and returns a Message pointer after having
// converted the v interface{} into a JSON string. The message is then
// indended to be sent to the Hub.Broadcast channel to be broadcast to the
// appropriate clients
func NewMessage(room, funcName string, v interface{}) (*Message, error) {
	body, err := toJSONString(v)
	if err != nil {
		return nil, err
	}
	return &Message{
		Room: room,
		Func: funcName,
		Body: body,
	}, nil
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Room: c.Room, Body: string(p)}
		c.Hub.Broadcast <- &message
		fmt.Printf("Message receieved: %+v\n", message)
	}
}
