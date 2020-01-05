package main

import socketio "github.com/googollee/go-socket.io"

func (app *application) setupSocketIO() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	app.socketIO = server

	return server, nil
}
