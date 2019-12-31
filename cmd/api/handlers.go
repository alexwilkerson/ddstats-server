package main

import (
	"net/http"
)

func (app *application) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Merle!"))
}
