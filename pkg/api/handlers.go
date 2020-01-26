package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/alexwilkerson/ddstats-server/pkg/websocket"
)

const (
	BronzeDaggerThreshold float64 = 60
	SilverDaggerThreshold float64 = 120
	GoldDaggerThreshold   float64 = 250
	DevilDaggerThreshold  float64 = 500
)

func (api *API) getDaily(w http.ResponseWriter, r *http.Request) {
	run, err := api.db.CollectorRuns.SelectMostRecent()
	if err != nil {
		api.serverError(w, err)
		return
	}
	newPlayers, err := api.db.CollectorNewPlayers.Select(run.ID)
	if err != nil {
		api.serverError(w, err)
		return
	}
	activePlayers, err := api.db.CollectorActivePlayers.Select(run.ID)
	if err != nil {
		api.serverError(w, err)
		return
	}
	highScores, err := api.db.CollectorHighScores.Select(run.ID)
	if err != nil {
		api.serverError(w, err)
		return
	}
	bronzeDaggers := []*models.CollectorHighScore{}
	silverDaggers := []*models.CollectorHighScore{}
	goldDaggers := []*models.CollectorHighScore{}
	devilDaggers := []*models.CollectorHighScore{}

	for _, player := range highScores {
		switch {
		case player.Score >= DevilDaggerThreshold:
			devilDaggers = append(devilDaggers, player)
		case player.Score >= GoldDaggerThreshold:
			goldDaggers = append(goldDaggers, player)
		case player.Score >= SilverDaggerThreshold:
			silverDaggers = append(silverDaggers, player)
		case player.Score >= BronzeDaggerThreshold:
			bronzeDaggers = append(bronzeDaggers, player)
		}
	}

	daily := struct {
		*models.CollectorRun
		NewPlayers    []*models.CollectorNewPlayer    `json:"new_players_list"`
		ActivePlayers []*models.CollectorActivePlayer `json:"active_players_list"`
		BronzeDaggers []*models.CollectorHighScore    `json:"bronze_daggers_list"`
		SilverDaggers []*models.CollectorHighScore    `json:"silver_daggers_list"`
		GoldDaggers   []*models.CollectorHighScore    `json:"gold_daggers_list"`
		DevilDaggers  []*models.CollectorHighScore    `json:"devil_daggers_list"`
	}{
		run,
		newPlayers,
		activePlayers,
		bronzeDaggers,
		silverDaggers,
		goldDaggers,
		devilDaggers,
	}

	api.writeJSON(w, daily)
}

func (api *API) getNews(w http.ResponseWriter, r *http.Request) {
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be an integer")
		return
	}
	if pageSize < 1 || pageSize > 100 {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be between 1 and 100")
		return
	}

	pageNum, err := strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be an integer")
		return
	}
	if pageNum < 1 {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be greater than 0")
		return
	}

	fmt.Println("pageSize", pageSize, "pageNum", pageNum)

	var news struct {
		TotalPages     int            `json:"total_pages"`
		TotalNewsCount int            `json:"total_news_count"`
		PageNumber     int            `json:"page_number"`
		PageSize       int            `json:"page_size"`
		NewsCount      int            `json:"news_count"`
		News           []*models.News `json:"news"`
	}

	news.News, err = api.db.News.GetAll(pageSize, pageNum)
	if err != nil {
		api.serverError(w, err)
		return
	}

	if news.News == nil {
		api.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	news.TotalNewsCount, err = api.db.News.GetTotalCount()
	if err != nil {
		api.serverError(w, err)
		return
	}

	news.TotalPages = int(math.Ceil(float64(news.TotalNewsCount) / float64(pageSize)))
	news.PageNumber = pageNum
	news.PageSize = pageSize
	news.NewsCount = len(news.News)

	api.writeJSON(w, news)
}

func (api *API) getReleases(w http.ResponseWriter, r *http.Request) {
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be an integer")
		return
	}
	if pageSize < 1 || pageSize > 100 {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be between 1 and 100")
		return
	}

	pageNum, err := strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be an integer")
		return
	}
	if pageNum < 1 {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be greater than 0")
		return
	}

	fmt.Println("pageSize", pageSize, "pageNum", pageNum)

	var releases struct {
		TotalPages        int               `json:"total_pages"`
		TotalReleaseCount int               `json:"total_releases_count"`
		PageNumber        int               `json:"page_number"`
		PageSize          int               `json:"page_size"`
		ReleaseCount      int               `json:"release_count"`
		Releases          []*models.Release `json:"releases"`
	}

	releases.Releases, err = api.db.Releases.GetAll(pageSize, pageNum)
	if err != nil {
		api.serverError(w, err)
		return
	}

	if releases.Releases == nil {
		api.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	releases.TotalReleaseCount, err = api.db.Releases.GetTotalCount()
	if err != nil {
		api.serverError(w, err)
		return
	}

	releases.TotalPages = int(math.Ceil(float64(releases.TotalReleaseCount) / float64(pageSize)))
	releases.PageNumber = pageNum
	releases.PageSize = pageSize
	releases.ReleaseCount = len(releases.Releases)

	api.writeJSON(w, releases)
}

func (api *API) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket endpoint hit")
	if _, ok := r.URL.Query()["room"]; !ok {
		api.clientError(w, http.StatusNotFound)
		return
	}

	room := r.URL.Query().Get("room")

	// upgrade this connection to a Websocket connection
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Hub:  api.websocketHub,
		Room: room,
	}

	api.websocketHub.Register <- client
	client.Read()
}

func (api *API) playerLive(w http.ResponseWriter, r *http.Request) {
	players := api.websocketHub.LivePlayers()
	api.writeJSON(w, struct {
		PlayerCount int                `json:"player_count"`
		Players     []websocket.Player `json:"players"`
	}{
		PlayerCount: len(players),
		Players:     players,
	})
}

func (api *API) submitGame(w http.ResponseWriter, r *http.Request) {
	var game models.SubmittedGame
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "malformed data")
		return
	}

	if game.PlayerID == -1 {
		api.clientMessage(w, http.StatusBadRequest, "some kind of error occurred")
		return
	}

	if game.Version == "" {
		api.clientMessage(w, http.StatusBadRequest, "ddstats version not found")
		return
	}

	if game.PlayerID == 0 {
		api.clientMessage(w, http.StatusBadRequest, "player ID not found")
		return
	}

	duplicate, id, err := api.db.SubmittedGames.CheckDuplicate(&game)
	if duplicate {
		api.writeJSON(w, struct {
			Message string `json:"message"`
			GameID  int    `json:"game_id"`
		}{"Replay already recorded.", id})
		return
	}
	if err != nil {
		api.serverError(w, err)
		return
	}

	// this retrieves the most recent player from the dd backend api and
	// updates the database this may take too much time, and if so...
	// it's worth it to take this block of code out and solely rely on the database.
	// it does, however ensure that each time a user submits a game, the user
	// data is up to date!
	player, err := api.ddAPI.UserByID(game.PlayerID)
	if err != nil {
		api.serverError(w, err)
		return
	}
	err = api.db.Players.UpsertDDPlayer(player)
	if err != nil {
		api.serverError(w, err)
		return
	}

	// This does the same as above, but for replay players.
	if game.ReplayPlayerID != 0 {
		replayPlayer, err := api.ddAPI.UserByID(game.ReplayPlayerID)
		if err != nil && !errors.Is(err, ddapi.ErrPlayerNotFound) {
			api.serverError(w, err)
			return
		}
		err = api.db.ReplayPlayers.Upsert(int(replayPlayer.PlayerID), replayPlayer.PlayerName)
		if err != nil {
			api.errorLog.Printf("%v", err)
		}
	}

	gameID, err := api.db.SubmittedGames.Insert(&game)
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "error while inserting data to database")
		return
	}

	api.writeJSON(w, struct {
		Message string `json:"message"`
		GameID  int    `json:"game_id"`
	}{"Game submitted.", gameID})
}

func (api *API) getGameFull(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	game, err := api.db.Games.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	states, err := api.db.Games.GetAll(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	v := struct {
		GameInfo *models.GameWithName `json:"game_info"`
		States   []*models.State      `json:"states"`
	}{
		GameInfo: game,
		States:   states,
	}

	api.writeJSON(w, v)
}

func (api *API) getGameAll(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetAll(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGameGems(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetGems(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGameHomingDaggers(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetHomingDaggers(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGameDaggersHit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetDaggersHit(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGameDaggersFired(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetDaggersFired(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGameAccuracy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetAccuracy(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGameEnemiesAlive(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetEnemiesAlive(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGameEnemiesKilled(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	states, err := api.db.Games.GetEnemiesKilled(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, states)
}

func (api *API) getGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	game, err := api.db.Games.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, game)
}

func (api *API) getPlayers(w http.ResponseWriter, r *http.Request) {
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be an integer")
		return
	}
	if pageSize < 1 || pageSize > 100 {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be between 1 and 100")
		return
	}

	pageNum, err := strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be an integer")
		return
	}
	if pageNum < 1 {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be greater than 0")
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

	players.Players, err = api.db.Players.GetAll(pageSize, pageNum)
	if err != nil {
		api.serverError(w, err)
		return
	}

	if players.Players == nil {
		api.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	players.TotalPlayerCount, err = api.db.Players.GetTotalCount()
	if err != nil {
		api.serverError(w, err)
		return
	}

	players.TotalPages = int(math.Ceil(float64(players.TotalPlayerCount) / float64(pageSize)))
	players.PageNumber = pageNum
	players.PageSize = pageSize
	players.PlayerCount = len(players.Players)

	api.writeJSON(w, players)
}

func (api *API) getRecentGames(w http.ResponseWriter, r *http.Request) {
	var playerID int
	var err error
	if _, ok := r.URL.Query()["player_id"]; ok {
		playerID, err = strconv.Atoi(r.URL.Query().Get("player_id"))
		if err != nil {
			api.clientMessage(w, http.StatusBadRequest, "player_id must be an integer")
			return
		}
		if playerID < 1 {
			api.clientMessage(w, http.StatusBadRequest, "player_id must be greater than 0")
			return
		}
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be an integer")
		return
	}
	if pageSize < 1 {
		api.clientMessage(w, http.StatusBadRequest, "page_size must be greater than 0")
		return
	}

	pageNum, err := strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be an integer")
		return
	}
	if pageNum < 1 {
		api.clientMessage(w, http.StatusBadRequest, "page_num must be greater than 0")
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

	games.Games, err = api.db.Games.GetRecent(playerID, pageSize, pageNum)
	if err != nil {
		api.serverError(w, err)
		return
	}

	if games.Games == nil {
		api.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	games.TotalGameCount, err = api.db.Games.GetTotalCount(playerID)
	if err != nil {
		api.serverError(w, err)
		return
	}

	games.TotalPages = int(math.Ceil(float64(games.TotalGameCount) / float64(pageSize)))
	games.PageNumber = pageNum
	games.PageSize = pageSize
	games.GameCount = len(games.Games)

	api.writeJSON(w, games)
}

func (api *API) getTopGames(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "limit must be an integer")
		return
	}
	if limit < 1 || limit > 100 {
		api.clientMessage(w, http.StatusBadRequest, "limit must be between 1 and 100")
		return
	}

	var games struct {
		GameCount int                    `json:"game_count"`
		Games     []*models.GameWithName `json:"games"`
	}

	games.Games, err = api.db.Games.GetTop(limit)
	if err != nil {
		api.serverError(w, err)
		return
	}

	if games.Games == nil {
		api.clientMessage(w, http.StatusNotFound, "no records found in this range")
		return
	}

	games.GameCount = len(games.Games)

	api.writeJSON(w, games)
}

func (api *API) getPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		api.clientError(w, http.StatusBadRequest)
		return
	}

	player, err := api.db.Players.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			api.clientMessage(w, http.StatusNotFound, err.Error())
		} else {
			api.serverError(w, err)
		}
		return
	}

	api.writeJSON(w, player)
}

func (api *API) playerUpdate(w http.ResponseWriter, r *http.Request) {
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

	err = api.db.Players.UpsertDDPlayer(player)
	if err != nil {
		api.clientMessage(w, http.StatusNotFound, "error updating player in database")
		fmt.Println(err)
		return
	}

	api.writeJSON(w, player)
}

func (api *API) getMOTD(w http.ResponseWriter, r *http.Request) {
	motd, err := api.db.MOTD.Get()
	if err != nil {
		api.serverError(w, err)
		return
	}

	api.writeJSON(w, motd)
}

func (api *API) clientConnect(w http.ResponseWriter, r *http.Request) {
	var version struct {
		Version string `json:"version"`
	}

	err := json.NewDecoder(r.Body).Decode(&version)
	if err != nil {
		api.clientMessage(w, http.StatusBadRequest, "malformed data")
		fmt.Println(err)
		return
	}

	motd, err := api.db.MOTD.Get()
	if err != nil {
		api.serverError(w, err)
		return
	}

	valid, err := validVersion(version.Version)
	if err != nil {
		api.serverError(w, err)
		return
	}
	update, err := updateAvailable(version.Version)
	if err != nil {
		api.serverError(w, err)
		return
	}

	data := struct {
		MOTD            string `json:"motd"`
		ValidVersion    bool   `json:"valid_version"`
		UpdateAvailable bool   `json:"update_available"`
	}{
		MOTD:            motd.Message,
		ValidVersion:    valid,
		UpdateAvailable: update,
	}

	api.writeJSON(w, data)
}
