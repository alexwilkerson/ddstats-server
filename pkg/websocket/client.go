package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   uint
	Conn *websocket.Conn
	Hub  *Hub
	Room string
}

type Message struct {
	Room string `json:"-"`
	Type int    `json:"type"`
	Func string `json:"func,omitempty"`
	Body string `json:"body,omitempty"`
}

func NewMessage(room, funcName string, v interface{}) (*Message, error) {
	body, err := toJSONString(v)
	if err != nil {
		return nil, err
	}
	return &Message{
		Room: room,
		Type: 1,
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
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Room: c.Room, Type: messageType, Body: string(p)}
		c.Hub.Broadcast <- &message
		fmt.Printf("Message receieved: %+v\n", message)
	}
}
