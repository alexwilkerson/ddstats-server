package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

func (app *application) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Merle!"))
}

func (app *application) showGame(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	game, err := app.games.GetGame(1)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)

		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Fprintf(w, "%v", game)
}
