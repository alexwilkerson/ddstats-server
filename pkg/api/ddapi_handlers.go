package api

import (
	"net/http"
	"strconv"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
)

func (api *API) ddGetUserByRank(w http.ResponseWriter, r *http.Request) {
	rank, err := strconv.Atoi(r.URL.Query().Get("rank"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "rank must be an integer")
		return
	}

	if rank < 1 {
		api.clientMessage(w, http.StatusBadRequest, "negative rank not allowed")
		return
	}

	player, err := api.ddAPI.UserByRank(rank)
	if err != nil {
		api.clientMessage(w, http.StatusNotFound, err.Error())
		return
	}

	api.writeJSON(w, player)
}

func (api *API) ddGetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	if id < 1 {
		api.clientMessage(w, http.StatusBadRequest, "negative id not allowed")
		return
	}

	player, err := api.ddAPI.UserByID(id)
	if err != nil {
		api.clientMessage(w, http.StatusNotFound, err.Error())
		return
	}

	api.writeJSON(w, player)
}

func (api *API) ddUserSearch(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		api.clientMessage(w, http.StatusBadRequest, "no 'name' query parameter set")
		return
	}
	api.infoLog.Println("name:", name)

	players, err := api.ddAPI.UserSearch(name)
	if err != nil {
		api.clientMessage(w, http.StatusNotFound, err.Error())
		return
	}

	data := struct {
		PlayerCount int             `json:"player_count"`
		Players     []*ddapi.Player `json:"players"`
	}{PlayerCount: len(players), Players: players}

	api.writeJSON(w, data)
}

func (api *API) ddGetScores(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "offset must be an integer")
		return
	}

	if offsetInt < 1 {
		api.clientMessage(w, http.StatusBadRequest, "offset must be greater than 0")
		return
	}

	limit := 100
	_, ok := r.URL.Query()["limit"]
	if ok {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			api.clientMessage(w, http.StatusBadRequest, "limit must be an integer")
			return
		}

		if limit < 1 || limit > 100 {
			api.clientMessage(w, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}
	}

	leaderboard, err := api.ddAPI.GetLeaderboard(limit, offsetInt)
	if err != nil {
		api.clientMessage(w, http.StatusNotFound, err.Error())
		return
	}

	api.writeJSON(w, leaderboard)
}
