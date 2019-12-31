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
	rank, ok := r.URL.Query()["rank"]
	if !ok || len(rank) < 1 {
		app.clientMessage(w, http.StatusBadRequest, "no 'rank' query parameter set")
		return
	}

	rankInt, err := strconv.Atoi(rank[0])
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "rank must be integer")
		return
	}

	if rankInt < 1 {
		app.clientMessage(w, http.StatusBadRequest, "negative rank not allowed")
		return
	}

	u := "http://dd.hasmodai.com/backend16/get_user_by_rank_public.php"
	form := url.Values{"rank": {rank[0]}}
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
	id, ok := r.URL.Query()["id"]
	if !ok || len(id) < 1 {
		app.clientMessage(w, http.StatusBadRequest, "no 'id' query parameter set")
		return
	}

	idInt, err := strconv.Atoi(id[0])
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "id must be integer")
		return
	}

	if idInt < 1 {
		app.clientMessage(w, http.StatusBadRequest, "negative id not allowed")
		return
	}

	u := "http://dd.hasmodai.com/backend16/get_user_by_id_public.php"
	form := url.Values{"uid": {id[0]}}
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
	name, ok := r.URL.Query()["name"]
	if !ok || len(name) < 1 {
		app.clientMessage(w, http.StatusBadRequest, "no 'name' query parameter set")
		return
	}

	u := "http://dd.hasmodai.com/backend16/get_user_search_public.php"
	form := url.Values{"search": {name[0]}}
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
