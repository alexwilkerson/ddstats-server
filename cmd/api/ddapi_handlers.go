package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var deathTypes = []string{
	"FALLEN", "SWARMED", "IMPALED", "GORED", "INFESTED", "OPENED", "PURGED",
	"DESECRATED", "SACRIFICED", "EVISCERATED", "ANNIHILATED", "INTOXICATED",
	"ENVENMONATED", "INCARNATED", "DISCARNATED", "BARBED",
}

// ErrPlayerNotFound returned when player not found from the DD API
var ErrPlayerNotFound = errors.New("player not found")

// Player is the struct returned after parsing the binary data
// blob returned from the DD API.
type Player struct {
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

	// start reading blob from byte position 19
	player, err := convertDDBytes(bodyBytes, 19)
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
	player, err := convertDDBytes(bodyBytes, 19)
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

func convertDDBytes(b []byte, bytePosition int) (*Player, error) {
	var player Player

	playerNameLength := int(toInt16(b, bytePosition))
	bytePosition += 2
	player.PlayerName = string(b[bytePosition : bytePosition+playerNameLength])
	bytePosition += playerNameLength
	// just figured out this information manually...
	player.PlayerID = toUint64(b, bytePosition+4)
	if player.PlayerID == 0 {
		return nil, ErrPlayerNotFound
	}
	player.Rank = toInt32(b, bytePosition)
	player.Time = float64(toInt32(b, bytePosition+12)) / 10000
	player.Kills = toInt32(b, bytePosition+16)
	player.Gems = toInt32(b, bytePosition+28)
	player.DaggersHit = toInt32(b, bytePosition+24)
	player.DaggersFired = toInt32(b, bytePosition+20)
	if player.DaggersFired > 0 {
		player.Accuracy = float64(player.DaggersHit) / float64(player.DaggersFired) * 100
	}
	player.DeathType = deathTypes[toInt16(b, bytePosition+32)]
	player.OverallTime = float64(toUint64(b, bytePosition+60)) / 10000
	player.OverallKills = toUint64(b, bytePosition+44)
	player.OverallGems = toUint64(b, bytePosition+68)
	player.OverallDeaths = toUint64(b, bytePosition+36)
	player.OverallDaggersHit = toUint64(b, bytePosition+76)
	player.OverallDaggersFired = toUint64(b, bytePosition+52)
	if player.OverallDaggersFired > 0 {
		player.OverallAccuracy = float64(player.OverallDaggersHit) / float64(player.OverallDaggersFired) * 100
	}

	return &player, nil
}

func toUint64(b []byte, offset int) uint64 {
	return binary.LittleEndian.Uint64(b[offset : offset+8])
}

func toInt64(b []byte, offset int) int64 {
	return int64(binary.LittleEndian.Uint64(b[offset : offset+4]))
}

func toUint32(b []byte, offset int) uint32 {
	return binary.LittleEndian.Uint32(b[offset : offset+4])
}

func toInt32(b []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(b[offset : offset+4]))
}

func toInt16(b []byte, offset int) int16 {
	return int16(binary.LittleEndian.Uint16(b[offset : offset+2]))
}
