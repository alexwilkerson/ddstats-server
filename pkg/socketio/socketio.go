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

	"github.com/alexwilkerson/ddstats-api/pkg/models/postgres"
	"github.com/alexwilkerson/ddstats-api/pkg/websocket"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	socketio "github.com/googollee/go-socket.io"
	"github.com/jmoiron/sqlx"
)

type sio struct {
	server       *socketio.Server
	client       *http.Client
	infoLog      *log.Logger
	errorLog     *log.Logger
	db           *sqlx.DB
	websocketHub *websocket.Hub
	ddAPI        *ddapi.API
	games        *postgres.GameModel
	players      *postgres.PlayerModel
	livePlayers  *sync.Map
}

const (
	defaultNamespace = "/"
)

type player struct {
	websocketPlayer *websocket.Player
	PlayerID        int     `json:"player_id"`
	PlayerName      string  `json:"player_name"`
	GameTime        float64 `json:"game_time"`
	DeathType       int     `json:"death_type"`
	IsReplay        bool    `json:"is_replay"`
}

type state struct {
	PlayerID         int     `json:"player_id"`
	GameTime         float64 `json:"game_time"`
	Gems             int     `json:"gems"`
	HomingDaggers    int     `json:"homing_daggers"`
	EnemiesAlive     int     `json:"enemies_alive"`
	EnemiesKilled    int     `json:"enemies_killed"`
	DaggersHit       int     `json:"daggers_hit"`
	DaggersFired     int     `json:"daggers_fired"`
	LevelTwoTime     float64 `json:"level_two_time"`
	LevelThreeTime   float64 `json:"level_three_time"`
	LevelFourTime    float64 `json:"level_four_time"`
	DeathType        int     `json:"death_type"`
	IsReplay         bool    `json:"is_replay"`
	NotifyPlayerBest bool    `json:"notify_player_best"`
	NotifyAbove1000  bool    `json:"notify_above_1000"`
}

// NewServer returns a Server from the go-socket.io package with all of the routes already
// set up to handle ddstats clients
func NewServer(infoLog, errorLog *log.Logger, websocketHub *websocket.Hub, client *http.Client, db *sqlx.DB) (*socketio.Server, error) {
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
		games:        &postgres.GameModel{DB: db},
		players:      &postgres.PlayerModel{DB: db},
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
}

// i don't know what this function should do
func (si *sio) onConnect(s socketio.Conn) error {
	s.SetContext("")
	si.infoLog.Println("connected:", s.ID())
	return nil
}

func (si *sio) onDisconnect(s socketio.Conn, msg string) {
	p, ok := si.livePlayers.Load(s.ID())
	if !ok {
		si.errorLog.Println("socketio onDisconnect: could not load player from livePlayers map")
		return
	}
	si.livePlayers.Delete(s.ID())
	si.websocketHub.UnregisterPlayer <- p.(*player).websocketPlayer
	si.infoLog.Println(s.ID(), "disconnected")
	return
}

func (si *sio) onLogin(s socketio.Conn, id int) {
	start := time.Now()
	// -1 is sent when there is an error in the client
	if id == -1 {
		si.errorLog.Println("socketio onLogin: id is -1")
		return
	}

	p, err := si.ddAPI.UserByID(id)
	if err != nil {
		si.errorLog.Printf("socketio onLogin: %w", err)
		return
	}

	websocketPlayer := websocket.Player{ID: int(p.PlayerID), Name: p.PlayerName}

	si.livePlayers.Store(s.ID(), &player{
		websocketPlayer: &websocketPlayer,
		PlayerID:        int(p.PlayerID),
		PlayerName:      p.PlayerName,
		GameTime:        0,
		DeathType:       -2, // IN MENU
		IsReplay:        false,
	})

	err = si.players.UpsertDDPlayer(p)
	if err != nil {
		si.errorLog.Printf("socketio onLogin: %w", err)
		return
	}

	si.websocketHub.RegisterPlayer <- &websocketPlayer

	si.infoLog.Println(id)
	si.infoLog.Println("duration:", time.Since(start))
}

func (si *sio) onSubmit(s socketio.Conn, playerID int, gameTime float64, gems, homingDaggers, enemiesAlive, enemiesKilled, daggersHit, daggersFired int, levelTwoTime, levelThreeTime, levelFourTime float64, isReplay bool, deathType int, notifyPlayerBest, notifyAbove1000 bool) {
	state := state{
		PlayerID:         playerID,
		GameTime:         gameTime,
		Gems:             gems,
		HomingDaggers:    homingDaggers,
		EnemiesAlive:     enemiesAlive,
		EnemiesKilled:    enemiesKilled,
		DaggersHit:       daggersHit,
		DaggersFired:     daggersFired,
		LevelTwoTime:     levelTwoTime,
		LevelThreeTime:   levelThreeTime,
		LevelFourTime:    levelFourTime,
		DeathType:        deathType,
		IsReplay:         isReplay,
		NotifyPlayerBest: notifyPlayerBest,
		NotifyAbove1000:  notifyAbove1000,
	}
	if playerID < 1 {
		return
	}
	websocketMessage, err := websocket.NewMessage(strconv.Itoa(playerID), "submit", state)
	if err != nil {
		si.errorLog.Println("socketio onSubmit: %w", err)
		return
	}
	si.websocketHub.Broadcast <- websocketMessage
	si.infoLog.Printf("%+v\n", state)
}

func (si *sio) onError(s socketio.Conn, err error) {
	si.errorLog.Printf("socketio onError: %w", err)
}
