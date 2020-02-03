// Package socketio runs the live stats on the website. it is intended to be backward compatible
// so that the client needn't be updated. However, it should be rewritten alongside
// the client in the future
package socketio

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
	"github.com/alexwilkerson/ddstats-server/pkg/websocket"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	socketio "github.com/googollee/go-socket.io"
)

const (
	StatusLoggedIn        = "Logged In"
	StatusNotConnected    = "Not Connected"
	StatusConnecting      = "Connecting"
	StatusAlive           = "Alive"
	StatusWatchingAReplay = "Watching A Replay"
	StatusInMainMenu      = "In Main Menu"
	StatusInDaggerLobby   = "In Dagger Lobby"
	StatusDead            = "Dead"
)

const (
	defaultNamespace = "/"
)

const (
	notifyThreshold = 1000
)

type sio struct {
	server       *socketio.Server
	client       *http.Client
	infoLog      *log.Logger
	errorLog     *log.Logger
	websocketHub *websocket.Hub
	ddAPI        *ddapi.API
	db           *postgres.Postgres
	livePlayers  *sync.Map
}

type player struct {
	sync.Mutex
	websocketPlayer        *websocket.PlayerWithLock
	PlayerID               int     `json:"player_id"`
	PlayerName             string  `json:"player_name"`
	BestGameTime           float64 `json:"best_game_time"`
	GameTime               float64 `json:"game_time"`
	DeathType              int     `json:"death_type"`
	IsReplay               bool    `json:"is_replay"`
	bestTimeNotified       bool
	aboveThresholdNotified bool
}

func (p *player) getStatus() string {
	var status string
	switch {
	case p.DeathType >= 0:
		status = StatusDead
	case p.DeathType == -2:
		status = StatusInMainMenu
	case p.DeathType == -1 && p.IsReplay == true:
		status = StatusWatchingAReplay
	default:
		status = StatusAlive
	}
	return status
}

type state struct {
	PlayerID             int     `json:"player_id"`
	GameTime             float64 `json:"game_time"`
	Gems                 int     `json:"gems"`
	HomingDaggers        int     `json:"homing_daggers"`
	EnemiesAlive         int     `json:"enemies_alive"`
	EnemiesKilled        int     `json:"enemies_killed"`
	DaggersHit           int     `json:"daggers_hit"`
	DaggersFired         int     `json:"daggers_fired"`
	LevelTwoTime         float64 `json:"level_two_time"`
	LevelThreeTime       float64 `json:"level_three_time"`
	LevelFourTime        float64 `json:"level_four_time"`
	DeathType            int     `json:"death_type"`
	IsReplay             bool    `json:"is_replay"`
	Status               string  `json:"status"`
	NotifyPlayerBest     bool    `json:"-"`
	NotifyAboveThreshold bool    `json:"-"`
}

// NewServer returns a Server from the go-socket.io package with all of the routes already
// set up to handle ddstats clients
func NewServer(infoLog, errorLog *log.Logger, websocketHub *websocket.Hub, client *http.Client, db *postgres.Postgres) (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	s := sio{
		server:       server,
		client:       client,
		infoLog:      infoLog,
		errorLog:     errorLog,
		db:           db,
		websocketHub: websocketHub,
		ddAPI:        &ddapi.API{Client: client},
		livePlayers:  &sync.Map{},
	}
	s.routes(server)
	return server, nil
}

func (si *sio) routes(server *socketio.Server) {
	server.OnConnect(defaultNamespace, si.onConnect)
	server.OnDisconnect(defaultNamespace, si.onDisconnect)
	server.OnError(defaultNamespace, si.onError)
	server.OnEvent(defaultNamespace, "login", si.onLogin)
	server.OnEvent(defaultNamespace, "submit", si.onSubmit)
	server.OnEvent(defaultNamespace, "status_update", si.onStatusUpdate)
	server.OnEvent(defaultNamespace, "game_submitted", si.onGameSubmitted)
}

// i don't know what this function should do
func (si *sio) onConnect(s socketio.Conn) error {
	s.SetContext("")
	si.infoLog.Println("connected:", s.ID())
	return nil
}

func (si *sio) onDisconnect(s socketio.Conn, msg string) {
	v, ok := si.livePlayers.Load(s.ID())
	if !ok {
		si.errorLog.Println("socketio onDisconnect: could not load player from livePlayers map")
		return
	}
	player := v.(*player)
	player.Lock()
	websocketMessage, err := websocket.NewMessage(strconv.Itoa(player.PlayerID), "submit", struct{}{})
	if err != nil {
		si.errorLog.Println("socketio onSubmit: %w", err)
	}
	si.websocketHub.Broadcast <- websocketMessage
	si.livePlayers.Delete(s.ID())
	si.websocketHub.UnregisterPlayer <- player.websocketPlayer
	player.Unlock()
	si.infoLog.Println(s.ID(), "disconnected")
	return
}

func (si *sio) onStatusUpdate(s socketio.Conn, playerID, statusID int) {
	var status string
	switch statusID {
	case 0:
		status = StatusNotConnected
	case 1:
		status = StatusConnecting
	case 2:
		status = StatusAlive
	case 3:
		status = StatusWatchingAReplay
	case 4:
		status = StatusInMainMenu
	case 5:
		status = StatusInDaggerLobby
	case 6:
		status = StatusDead
	}
	p, ok := si.livePlayers.Load(s.ID())
	if !ok {
		si.errorLog.Printf("player with s.ID() %s not found in livePlayers map", s.ID())
		return
	}
	player := p.(*player)
	player.Lock()
	player.websocketPlayer.Lock()
	player.websocketPlayer.Status = status
	player.websocketPlayer.Unlock()
	player.Unlock()

	websocketMessage, err := websocket.NewMessage(strconv.Itoa(playerID), "status", status)
	if err != nil {
		si.errorLog.Println("socketio status: %w", err)
	}
	si.websocketHub.Broadcast <- websocketMessage
}

func (si *sio) onGameSubmitted(s socketio.Conn, gameID int, notifyPlayerBest, notifyAboveThreshold bool) {
	v, ok := si.livePlayers.Load(s.ID())
	if !ok {
		return
	}
	player := v.(*player)
	player.Lock()
	defer player.Unlock()
	game, err := si.db.Games.Get(gameID)
	if err != nil {
		si.errorLog.Printf("%+v", err)
	}
	if notifyPlayerBest && game.GameTime > player.BestGameTime {
		si.websocketHub.DiscordBroadcast <- &websocket.PlayerBestSubmitted{
			PlayerName:       player.PlayerName,
			GameID:           gameID,
			GameTime:         game.GameTime,
			PreviousGameTime: player.BestGameTime,
		}
	}
	if notifyAboveThreshold && game.GameTime > notifyThreshold {
		si.websocketHub.DiscordBroadcast <- &websocket.PlayerDied{
			PlayerName: player.PlayerName,
			GameID:     gameID,
			GameTime:   game.GameTime,
			DeathType:  game.DeathType,
		}
	}
}

func (si *sio) onLogin(s socketio.Conn, id int) {
	start := time.Now()
	// -1 is sent when there is an error in the client
	if id == -1 {
		si.errorLog.Println("socketio onLogin: id is -1")
		s.Close()
		return
	}

	p, err := si.ddAPI.UserByID(id)
	if err != nil {
		si.errorLog.Printf("socketio onLogin: %w", err)
		s.Close()
		return
	}

	websocketPlayer := websocket.PlayerWithLock{Player: websocket.Player{ID: int(p.PlayerID), Name: p.PlayerName, Status: StatusLoggedIn}}

	si.livePlayers.Store(s.ID(), &player{
		websocketPlayer: &websocketPlayer,
		PlayerID:        int(p.PlayerID),
		PlayerName:      p.PlayerName,
		BestGameTime:    p.GameTime,
		DeathType:       -2, // IN MENU
	})

	err = si.db.Players.UpsertDDPlayer(p)
	if err != nil {
		si.errorLog.Printf("socketio onLogin: %w", err)
		s.Close()
		return
	}

	si.websocketHub.RegisterPlayer <- &websocketPlayer

	websocketMessage, err := websocket.NewMessage(strconv.Itoa(int(p.PlayerID)), "submit", struct{}{})
	if err != nil {
		si.errorLog.Println("socketio onSubmit: %w", err)
	}
	si.websocketHub.Broadcast <- websocketMessage

	si.infoLog.Println(id)
	si.infoLog.Println("duration:", time.Since(start))
}

func (si *sio) onSubmit(s socketio.Conn, playerID int, gameTime float64, gems, homingDaggers, enemiesAlive, enemiesKilled, daggersHit, daggersFired int, levelTwoTime, levelThreeTime, levelFourTime float64, isReplay bool, deathType int, notifyPlayerBest, notifyAboveThreshold bool) {
	state := state{
		PlayerID:             playerID,
		GameTime:             gameTime,
		Gems:                 gems,
		HomingDaggers:        homingDaggers,
		EnemiesAlive:         enemiesAlive,
		EnemiesKilled:        enemiesKilled,
		DaggersHit:           daggersHit,
		DaggersFired:         daggersFired,
		LevelTwoTime:         levelTwoTime,
		LevelThreeTime:       levelThreeTime,
		LevelFourTime:        levelFourTime,
		DeathType:            deathType,
		IsReplay:             isReplay,
		NotifyPlayerBest:     notifyPlayerBest,
		NotifyAboveThreshold: notifyAboveThreshold,
	}
	if playerID < 1 {
		si.errorLog.Println("playerID less than 1")
		return
	}
	p, ok := si.livePlayers.Load(s.ID())
	if !ok {
		si.errorLog.Printf("player with s.ID() %s not found in livePlayers map", s.ID())
		return
	}
	player := p.(*player)
	player.Lock()
	state.Status = player.getStatus()
	defer player.Unlock()
	if state.GameTime < player.GameTime {
		player.aboveThresholdNotified = false
		player.bestTimeNotified = false
	}
	player.GameTime = state.GameTime
	player.DeathType = state.DeathType
	player.IsReplay = state.IsReplay
	status := player.getStatus()
	player.websocketPlayer.Lock()
	player.websocketPlayer.GameTime = state.GameTime
	player.websocketPlayer.Status = status
	player.websocketPlayer.Unlock()
	websocketMessage, err := websocket.NewMessage(strconv.Itoa(playerID), "submit", state)
	if err != nil {
		si.errorLog.Println("socketio onSubmit: %w", err)
		return
	}
	si.websocketHub.Broadcast <- websocketMessage
	// if the notification hasn't yet happened and a player beats their previous score,
	if notifyPlayerBest && !player.bestTimeNotified && gameTime > player.BestGameTime {
		si.websocketHub.DiscordBroadcast <- &websocket.PlayerBestReached{
			PlayerID:         player.PlayerID,
			PlayerName:       player.PlayerName,
			PreviousGameTime: player.BestGameTime,
		}
		player.bestTimeNotified = true
	}
	if notifyAboveThreshold && !player.aboveThresholdNotified && gameTime >= notifyThreshold {
		player.aboveThresholdNotified = true
		si.websocketHub.DiscordBroadcast <- &websocket.PlayerAboveThreshold{
			PlayerID:   player.PlayerID,
			PlayerName: player.PlayerName,
		}
	}
}

func (si *sio) onError(s socketio.Conn, err error) {
	si.errorLog.Printf("socketio onError: %+v", err)
}
