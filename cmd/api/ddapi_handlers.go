package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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

	u := "http://dd.hasmodai.com/backend16/get_user_by_rank_public.php"
	form := url.Values{"rank": {strconv.Itoa(rank)}}
	resp, err := app.client.PostForm(u, form)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		app.serverError(w, err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// start reading blob from byte position 19
	player, err := ddapi.BytesToPlayer(bodyBytes, 19)
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

	u := "http://dd.hasmodai.com/backend16/get_user_by_id_public.php"
	form := url.Values{"uid": {strconv.Itoa(id)}}
	resp, err := app.client.PostForm(u, form)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		app.serverError(w, err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// start reading blob from byte position 19
	player, err := ddapi.BytesToPlayer(bodyBytes, 19)
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

	u := "http://dd.hasmodai.com/backend16/get_user_search_public.php"
	form := url.Values{"search": {name}}
	resp, err := app.client.PostForm(u, form)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		app.serverError(w, err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	players, err := ddapi.UserSearchBytesToPlayers(bodyBytes)
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
