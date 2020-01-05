package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

func (app *application) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Amer!"))
}

func (app *application) submitGame(w http.ResponseWriter, r *http.Request) {
	var game models.SubmittedGame
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "malformed data")
		return
	}

	if game.PlayerID == -1 {
		app.clientMessage(w, http.StatusBadRequest, "some kind of error occurred")
		return
	}

	if game.Version == "" {
		app.clientMessage(w, http.StatusBadRequest, "ddstats version not found")
		return
	}

	if game.PlayerID == 0 {
		app.clientMessage(w, http.StatusBadRequest, "player ID not found")
		return
	}

	gameID, err := app.submittedGames.Insert(&game)
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "error while inserting data to database")
		return
	}

	js, err := json.Marshal(struct {
		Message string `json:"message"`
		GameID  int    `json:"game_id"`
	}{"Game submitted.", gameID})
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "error retrieving game ID")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

func (app *application) getGameGems(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetGems(id)
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

func (app *application) getGameHomingDaggers(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetHomingDaggers(id)
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

func (app *application) getGameDaggersHit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetDaggersHit(id)
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

func (app *application) getGameDaggersFired(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetDaggersFired(id)
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

func (app *application) getGameAccuracy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetAccuracy(id)
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

func (app *application) getGameEnemiesAlive(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetEnemiesAlive(id)
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

func (app *application) getGameEnemiesKilled(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := app.games.GetEnemiesKilled(id)
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
	if pageSize < 1 || pageSize > 100 {
		app.clientMessage(w, http.StatusBadRequest, "pagesize must be between 1 and 100")
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

	players.TotalPlayerCount, err = app.players.GetTotalCount()
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

func (app *application) getRecentGames(w http.ResponseWriter, r *http.Request) {
	var playerID int
	var err error
	if _, ok := r.URL.Query()["playerid"]; ok {
		playerID, err = strconv.Atoi(r.URL.Query().Get("playerid"))
		if err != nil {
			app.clientMessage(w, http.StatusBadRequest, "playerid must be an integer")
			return
		}
		if playerID < 1 {
			app.clientMessage(w, http.StatusBadRequest, "playerid must be greater than 0")
			return
		}
	}

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

	var games struct {
		TotalPages     int                    `json:"total_pages"`
		TotalGameCount int                    `json:"total_game_count"`
		PageNumber     int                    `json:"page_number"`
		PageSize       int                    `json:"page_size"`
		GameCount      int                    `json:"game_count"`
		Games          []*models.GameWithName `json:"games"`
	}

	games.Games, err = app.games.GetRecent(playerID, pageSize, pageNum)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if games.Games == nil {
		app.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	games.TotalGameCount, err = app.games.GetTotalCount(playerID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	games.TotalPages = int(math.Ceil(float64(games.TotalGameCount) / float64(pageSize)))
	games.PageNumber = pageNum
	games.PageSize = pageSize
	games.GameCount = len(games.Games)

	js, err := json.Marshal(games)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) getTopGames(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "limit must be an integer")
		return
	}
	if limit < 1 || limit > 100 {
		app.clientMessage(w, http.StatusBadRequest, "limit must be between 1 and 100")
		return
	}

	var games struct {
		GameCount int                    `json:"game_count"`
		Games     []*models.GameWithName `json:"games"`
	}

	games.Games, err = app.games.GetTop(limit)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if games.Games == nil {
		app.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	games.GameCount = len(games.Games)

	js, err := json.Marshal(games)
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

func (app *application) playerUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.clientMessage(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	if id < 1 {
		app.clientMessage(w, http.StatusBadRequest, "negative id not allowed")
		return
	}

	player, err := app.ddAPI.UserByID(id)
	if err != nil {
		app.clientMessage(w, http.StatusNotFound, err.Error())
		return
	}

	err = app.players.UpsertDDPlayer(player)
	if err != nil {
		app.clientMessage(w, http.StatusNotFound, "error updating player in database")
		fmt.Println(err)
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

func (app *application) getMOTD(w http.ResponseWriter, r *http.Request) {
	motd, err := app.motd.Get()
	if err != nil {
		app.serverError(w, err)
		return
	}

	js, err := json.Marshal(motd)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
