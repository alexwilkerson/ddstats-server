package models

import (
	"database/sql"
	"errors"
	"time"
)

//ErrNoRecord will be returned when DB record not found
var ErrNoRecord = errors.New("no record found")

//Game record representation
type Game struct {
	ID                   uint           `json:"id"`
	PlayerID             uint           `json:"player_id"`
	Granularity          int            `json:"granularity"`
	GameTime             float64        `json:"game_time"`
	DeathType            int            `json:"death_type"`
	Gems                 uint           `json:"gems"`
	HomingDaggers        uint           `json:"homing_daggers"`
	DaggersFired         uint           `json:"daggers_fired"`
	DaggersHit           uint           `json:"daggers_hit"`
	EnemiesAlive         uint           `json:"enemies_alive"`
	EnemiesKilled        uint           `json:"enemies_killed"`
	TimeStamp            time.Time      `json:"time_stamp"`
	ReplayPlayerID       int            `json:"replay_player_id"`
	SurvivalHash         string         `json:"survival_hash"`
	Version              sql.NullString `json:"version"`
	LevelTwoTime         float64        `json:"level_two_time"`
	LevelThreeTime       float64        `json:"level_three_time"`
	LevelFourTime        float64        `json:"level_four_time"`
	HomingDaggersMaxTime float64        `json:"homing_daggers_max_time"`
	EnemiesAliveMaxTime  float64        `json:"enemies_alive_max_time"`
	HomingDaggersMax     uint           `json:"homing_daggers_max"`
	EnemiesAliveMax      uint           `json:"enemies_alive_max"`
}
