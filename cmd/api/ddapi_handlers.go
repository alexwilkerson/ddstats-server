package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
)

func (app *application) ddGetUserByRank(w http.ResponseWriter, r *http.Request) {
	rank, err := strconv.Atoi(r.URL.Query().Get("rank"))
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "rank must be an integer")
		return
	}

	if rank < 1 {
		app.clientMessage(w, http.StatusBadRequest, "negative rank not allowed")
		return
	}

	player, err := app.ddAPI.UserByRank(rank)
	if err != nil {
		app.clientMessage(w, http.StatusNotFound, err.Error())
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

func (app *application) ddGetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	if id < 1 {
		app.clientMessage(w, http.StatusBadRequest, "negative id not allowed")
		return
	}

	// start reading blob from byte position 19
	player, err := app.ddAPI.UserByID(id)
	if err != nil {
		app.clientMessage(w, http.StatusNotFound, err.Error())
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

func (app *application) ddUserSearch(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		app.clientMessage(w, http.StatusBadRequest, "no 'name' query parameter set")
		return
	}

	players, err := app.ddAPI.UserSearch(name)
	if err != nil {
		app.clientMessage(w, http.StatusNotFound, err.Error())
		return
	}

	data := struct {
		PlayerCount int             `json:"player_count"`
		Players     []*ddapi.Player `json:"players"`
	}{PlayerCount: len(players), Players: players}

	js, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) ddGetScores(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "offset must be an integer")
		return
	}

	if offsetInt < 1 {
		app.clientMessage(w, http.StatusBadRequest, "offset must be greater than 0")
		return
	}

	limit := 100
	_, ok := r.URL.Query()["limit"]
	if ok {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			app.clientMessage(w, http.StatusBadRequest, "limit must be an integer")
			return
		}

		if limit < 1 || limit > 100 {
			app.clientMessage(w, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}
	}

	leaderboard, err := app.ddAPI.GetLeaderboard(limit, offsetInt)
	if err != nil {
		app.clientMessage(w, http.StatusNotFound, err.Error())
		return
	}

	js, err := json.Marshal(leaderboard)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
