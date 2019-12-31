package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (app *application) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Merle!"))
}

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

	data := struct {
		PlayerName          string  `json:"player_name"`
		PlayerID            uint64  `json:"player_id"`
		Rank                int32   `json:"rank"`
		Time                float64 `json:"time"`
		Kills               int32   `json:"kills"`
		Gems                int32   `json:"gems"`
		DaggersHit          int32   `json:"daggers_hit"`
		DaggersFired        int32   `json:"daggers_fired"`
		Accuracy            float64 `json:"accuracy"`
		DeathType           string  `json:"death_type"`
		OverallTime         float64 `json:"overall_time"`
		OverallKills        uint64  `json:"overall_kills"`
		OverallGems         uint64  `json:"overall_gems"`
		OverallDeaths       uint64  `json:"overall_deaths"`
		OverallDaggersHit   uint64  `json:"overall_daggers_hit"`
		OverallDaggersFired uint64  `json:"overall_daggers_fired"`
		OverallAccuracy     float64 `json:"overall_accuracy"`
	}{}

	// start reading blob from here
	offset := 19

	playerNameLength := int(toInt16(bodyBytes, offset))
	offset += 2
	data.PlayerName = string(bodyBytes[offset : offset+playerNameLength])
	offset += playerNameLength
	// just figured out this information manually...
	data.PlayerID = toUint64(bodyBytes, offset+4)
	if data.PlayerID == 0 {
		app.clientMessage(w, http.StatusNotFound, "no player found")
		return
	}
	data.Rank = toInt32(bodyBytes, offset)
	data.Time = float64(toInt32(bodyBytes, offset+12)) / 10000
	data.Kills = toInt32(bodyBytes, offset+16)
	data.Gems = toInt32(bodyBytes, offset+28)
	data.DaggersHit = toInt32(bodyBytes, offset+24)
	data.DaggersFired = toInt32(bodyBytes, offset+20)
	if data.DaggersFired > 0 {
		data.Accuracy = float64(data.DaggersHit) / float64(data.DaggersFired) * 100
	}
	data.DeathType = deathTypes[toInt16(bodyBytes, offset+32)]
	data.OverallTime = float64(toUint64(bodyBytes, offset+60)) / 10000
	data.OverallKills = toUint64(bodyBytes, offset+44)
	data.OverallGems = toUint64(bodyBytes, offset+68)
	data.OverallDeaths = toUint64(bodyBytes, offset+36)
	data.OverallDaggersHit = toUint64(bodyBytes, offset+76)
	data.OverallDaggersFired = toUint64(bodyBytes, offset+52)
	if data.OverallDaggersFired > 0 {
		data.OverallAccuracy = float64(data.OverallDaggersHit) / float64(data.OverallDaggersFired) * 100
	}

	js, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
