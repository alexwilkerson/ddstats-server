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

	if offsetInt < 0 {
		app.clientMessage(w, http.StatusBadRequest, "negative offset not allowed")
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

	// the DD API weirdly counts users starting from 1 but internally uses a 0 index
	// this fix it to make it more readable for users.
	if offsetInt != 0 {
		offsetInt--
	}

	u := "http://dd.hasmodai.com/backend16/get_scores.php"
	form := url.Values{"user": {"0"}, "level": {"survival"}, "offset": {strconv.Itoa(offsetInt)}}
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

	leaderboard, err := ddapi.GetScoresBytesToLeaderboard(bodyBytes, limit)
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
