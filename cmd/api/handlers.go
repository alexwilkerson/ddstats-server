package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

func (app *application) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Merle!"))
}

func (app *application) showGame(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.clientMessage(w, 400, "Bad query params request")
		return
	}

	game, err := app.games.GetGame(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientMessage(w, 404, "Game Does not exist")

		} else {
			app.serverError(w, err)
		}
		return
	}
	//This should work since data had been retreived
	gameValue, err := json.Marshal(game)
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(gameValue)

}
