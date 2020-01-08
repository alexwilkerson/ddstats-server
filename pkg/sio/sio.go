// Package sio runs the live stats on the website. it is intended to be backward compatible
// so that the client needn't be updated. However, it should be rewritten alongside
// the client in the future
package sio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/models/postgres"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	socketio "github.com/googollee/go-socket.io"
	"github.com/jmoiron/sqlx"
)

type sio struct {
	server      *socketio.Server
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
	PlayerID   int `json:"player_id"`
	PlayerName string
	GameTime   float64
	DeathType  int
	IsReplay   bool
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
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	s := sio{
		server:      server,
		client:      client,
		db:          db,
		ddAPI:       &ddapi.API{Client: client},
		games:       &postgres.GameModel{DB: db},
		players:     &postgres.PlayerModel{DB: db},
		livePlayers: map[string]*player{},
	}
	s.routes(server)
	return server, nil
}

func (si *sio) routes(server *socketio.Server) {
	server.OnConnect(defaultNamespace, si.onConnect)
	server.OnDisconnect(defaultNamespace, si.onDisconnect)
	server.OnEvent(defaultNamespace, "login", si.onLogin)
	server.OnEvent(defaultNamespace, "submit", si.onSubmit)
	server.OnEvent(defaultNamespace, "hello", func(s socketio.Conn, msg string) {
		fmt.Println(msg)
	})
	// server.OnEvent(defaultNamespace, "hello", func(s socketio.Conn, msg string) {
	// 	fmt.Println("helloed")
	// 	fmt.Println(msg)
	// })
}

func (si *sio) onConnect(s socketio.Conn) error {
	s.SetContext("")
	p := player{
		151515,
		"vhs",
		25.02,
		135,
		true,
	}
	si.livePlayers[s.ID()] = &p
	js, err := json.Marshal(si.livePlayers)
	if err != nil {
		return err
	}
	fmt.Println(string(js))
	s.Emit("live_users_update", string(js))
	// s.Emit("live_users_update", "oahtoahwtoh")

	fmt.Println("connected:", s.ID())
	return nil
}

func (si *sio) onDisconnect(s socketio.Conn, msg string) {
	if _, ok := si.livePlayers[s.ID()]; !ok {
		fmt.Println("this")
		return
	}
	player := si.livePlayers[s.ID()]
	delete(si.livePlayers, s.ID())
	js, err := json.Marshal(si.livePlayers)
	if err != nil {
		fmt.Println("that")
		return
	}
	fmt.Println(s.ID(), "disconnected")
	s.Emit("live_users_update", string(js))
	si.server.BroadcastToRoom(strconv.Itoa(player.PlayerID), "offline")
	return
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
		PlayerID:   int(p.PlayerID),
		PlayerName: p.PlayerName,
		GameTime:   0,
		DeathType:  -2, // IN MENU
		IsReplay:   false,
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
