package models

import (
	"database/sql"
	"errors"
	"time"
)

//ErrNoRecord will be returned when DB record not found
var ErrNoRecord = errors.New("No record found")

//Game reccord representation
type Game struct {
	ID                   uint
	PlayerID             uint
	Granularity          int
	GameTime             float64
	DeathType            int
	Gems                 uint
	HomingDaggers        uint
	DaggersFired         uint
	DaggersHit           uint
	EnemiesAlive         uint
	EnemiesKilled        uint
	TimeStamp            time.Time
	ReplayPlayerID       int
	SurvivalHash         string
	Version              sql.NullString
	LevelTwoTime         float64
	LevelThreeTime       float64
	LevelFourTime        float64
	HomingDaggersMaxTime float64
	EnemiesAliveMaxTime  float64
	HomingDaggersMax     uint
	EnemiesAliveMax      uint
}
