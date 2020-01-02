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

func (app *application) getGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	game, err := app.games.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			app.serverError(w, err)
		}
		return
	}

	js, err := json.Marshal(game)
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) getPlayers(w http.ResponseWriter, r *http.Request) {
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pagesize"))
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "pagesize must be an integer")
		return
	}
	if pageSize < 1 {
		app.clientMessage(w, http.StatusBadRequest, "pagesize must be greater than 0")
		return
	}

	pageNum, err := strconv.Atoi(r.URL.Query().Get("pagenum"))
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "pagenum must be an integer")
		return
	}
	if pageNum < 1 {
		app.clientMessage(w, http.StatusBadRequest, "pagenum must be greater than 0")
		return
	}

	var players struct {
		Players []*models.Player `json:"players"`
	}

	players.Players, err = app.players.GetAll(pageSize, pageNum)
	if err != nil {
		app.serverError(w, err)
		return
	}

	js, err := json.Marshal(players)
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) getPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	player, err := app.players.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			app.serverError(w, err)
		}
		return
	}

	js, err := json.Marshal(player)
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
