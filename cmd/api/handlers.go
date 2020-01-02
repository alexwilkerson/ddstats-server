package main

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

func (app *application) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Merle!"))
}

func (app *application) getGameAll(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetAll(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			app.serverError(w, err)
		}
		return
	}

	js, err := json.Marshal(states)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
		app.serverError(w, err)
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
		TotalPages       int              `json:"total_pages"`
		TotalPlayerCount int              `json:"total_player_count"`
		PageNumber       int              `json:"page_number"`
		PageSize         int              `json:"page_size"`
		PlayerCount      int              `json:"player_count"`
		Players          []*models.Player `json:"players"`
	}

	players.Players, err = app.players.GetAll(pageSize, pageNum)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if players.Players == nil {
		app.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	players.TotalPlayerCount, err = app.players.GetPlayerCount()
	if err != nil {
		app.serverError(w, err)
		return
	}

	players.TotalPages = int(math.Ceil(float64(players.TotalPlayerCount) / float64(pageSize)))
	players.PageNumber = pageNum
	players.PageSize = pageSize
	players.PlayerCount = len(players.Players)

	js, err := json.Marshal(players)
	if err != nil {
		app.serverError(w, err)
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
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
