// Package sio runs the live stats on the website. it is intended to be backward compatible
// so that the client needn't be updated. However, it should be rewritten alongside
// the client in the future
package sio

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/models/postgres"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	socketio "github.com/googollee/go-socket.io"
	"github.com/jmoiron/sqlx"
)

type sio struct {
	client      *http.Client
	db          *sqlx.DB
	ddAPI       *ddapi.API
	games       *postgres.GameModel
	players     *postgres.PlayerModel
	livePlayers map[string]*player
}

const (
	defaultNamespace  = "/"
	botNamespace      = "/ddstats-bot"
	userPageNamespace = "/user_page"
	indexNamespace    = "/index"
)

type player struct {
	playerID   int
	playerName string
	gameTime   float64
	deathType  int
	isReplay   bool
}

type state struct {
	PlayerID         int
	GameTime         float64
	Gems             int
	HomingDaggers    int
	EnemiesAlive     int
	EnemiesKilled    int
	DaggersHit       int
	DaggersFired     int
	LevelTwoTime     float64
	LevelThreeTime   float64
	LevelFourTime    float64
	DeathType        int
	IsReplay         bool
	NotifyPlayerBest bool
	NotifyAbove1000  bool
}

func Server(client *http.Client, db *sqlx.DB) (*socketio.Server, error) {
	s := sio{
		client:      client,
		db:          db,
		ddAPI:       &ddapi.API{Client: client},
		games:       &postgres.GameModel{DB: db},
		players:     &postgres.PlayerModel{DB: db},
		livePlayers: map[string]*player{},
	}
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	s.routes(server)
	return server, nil
}

func (si *sio) routes(server *socketio.Server) {
	server.OnConnect(defaultNamespace, si.onConnect)
	server.OnEvent(defaultNamespace, "login", si.onLogin)
	server.OnEvent(defaultNamespace, "submit", si.onSubmit)
}

func (si *sio) onConnect(s socketio.Conn) error {
	s.SetContext("")

	s.Emit("live_users_update", "working")

	fmt.Println("connected:", s.ID())
	return nil
}

func (si *sio) onLogin(s socketio.Conn, id int) {
	start := time.Now()
	// -1 is sent when there is an error in the client
	if id == -1 {
		// todo: handle error, print?
		return
	}

	p, err := si.ddAPI.UserByID(id)
	if err != nil {
		// todo: handle error, print?
		return
	}

	si.livePlayers[s.ID()] = &player{
		playerID:   int(p.PlayerID),
		playerName: p.PlayerName,
		gameTime:   0,
		deathType:  -2, // IN MENU
		isReplay:   false,
	}

	err = si.players.UpsertDDPlayer(p)
	if err != nil {
		// todo: handle error, print?
		return
	}

	fmt.Println(id)
	fmt.Println("duration:", time.Since(start))
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
	if playerID == -1 {
		return
	}
	fmt.Printf("%+v\n", state)
}
