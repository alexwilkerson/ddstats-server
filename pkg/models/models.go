package models

import (
	"errors"
	"time"

	"gopkg.in/guregu/null.v3"
)

//ErrNoRecord will be returned when DB record not found
var ErrNoRecord = errors.New("no record found")

//Game record representation
type Game struct {
	ID                   int         `json:"id"`
	PlayerID             int         `json:"player_id"`
	Granularity          int         `json:"granularity"`
	GameTime             float64     `json:"game_time"`
	DeathType            int         `json:"death_type"`
	Gems                 int         `json:"gems"`
	HomingDaggers        int         `json:"homing_daggers"`
	DaggersFired         int         `json:"daggers_fired"`
	DaggersHit           int         `json:"daggers_hit"`
	Accuracy             float64     `json:"accuracy"`
	EnemiesAlive         int         `json:"enemies_alive"`
	EnemiesKilled        int         `json:"enemies_killed"`
	TimeStamp            time.Time   `json:"time_stamp"`
	ReplayPlayerID       int         `json:"replay_player_id"`
	SurvivalHash         string      `json:"survival_hash"`
	Version              null.String `json:"version"`
	LevelTwoTime         float64     `json:"level_two_time"`
	LevelThreeTime       float64     `json:"level_three_time"`
	LevelFourTime        float64     `json:"level_four_time"`
	HomingDaggersMaxTime float64     `json:"homing_daggers_max_time"`
	EnemiesAliveMaxTime  float64     `json:"enemies_alive_max_time"`
	HomingDaggersMax     int         `json:"homing_daggers_max"`
	EnemiesAliveMax      int         `json:"enemies_alive_max"`
}

// Player struct is for players
type Player struct {
	ID                   int     `json:"id"`
	PlayerName           string  `json:"player_name"`
	Rank                 int     `json:"rank"`
	GameTime             float64 `json:"game_time"`
	DeathType            string  `json:"death_type"`
	Gems                 int     `json:"gems"`
	DaggersHit           int     `json:"daggers_hit"`
	DaggersFired         int     `json:"daggers_fired"`
	EnemiesKilled        int     `json:"enemies_killed"`
	Accuracy             float64 `json:"accuracy"`
	OverallTime          float64 `json:"overall_time"`
	OverallDeaths        int     `json:"overall_deaths"`
	OverallGems          int     `json:"overall_gems"`
	OverallEnemiesKilled int     `json:"overall_enemies_killed"`
	OverallDaggersHit    int     `json:"overall_daggers_hit"`
	OverallDaggersFired  int     `json:"overall_daggers_fired"`
	OverallAccuracy      float64 `json:"overall_accuracy"`
}

// State struct is for State
type State struct {
	GameID        int     `json:"game_id,omitempty"`
	GameTime      float64 `json:"game_time"`
	Gems          int     `json:"gems"`
	HomingDaggers int     `json:"homing_daggers"`
	DaggersHit    int     `json:"daggers_hit"`
	DaggersFired  int     `json:"daggers_fired"`
	Accuracy      float64 `json:"accuracy"`
	EnemiesAlive  int     `json:"enemies_alive"`
	EnemiesKilled int     `json:"enemies_killed"`
}

type Gems struct {
	GameTime float64 `json:"game_time"`
	Gems     int     `json:"gems"`
}

type HomingDaggers struct {
	GameTime      float64 `json:"game_time"`
	HomingDaggers int     `json:"homing_daggers"`
}

type DaggersHit struct {
	GameTime   float64 `json:"game_time"`
	DaggersHit int     `json:"daggers_hit"`
}

type DaggersFired struct {
	GameTime     float64 `json:"game_time"`
	DaggersFired int     `json:"daggers_fired"`
}

type Accuracy struct {
	GameTime float64 `json:"game_time"`
	Accuracy float64 `json:"accuracy"`
}

type EnemiesAlive struct {
	GameTime     float64 `json:"game_time"`
	EnemiesAlive int     `json:"enemies_alive"`
}

type EnemiesKilled struct {
	GameTime      float64 `json:"game_time"`
	EnemiesKilled int     `json:"enemies_killed"`
}

type SubmittedGame struct {
	PlayerID            int       `json:"playerID"`
	PlayerName          string    `json:"playerName"`
	Granularity         int       `json:"granularity"`
	GameTime            float64   `json:"inGameTimer"`
	GameTimeSlice       []float64 `json:"inGameTimerVector"`
	Gems                int       `json:"gems"`
	GemsSlice           []int     `json:"gemsVector"`
	Level2Time          float64   `json:"levelTwoTime"`
	Level3Time          float64   `json:"levelThreeTime"`
	Level4Time          float64   `json:"levelFourTime"`
	HomingDaggers       int       `json:"homingDaggers"`
	HomingDaggersSlice  []int     `json:"homingDaggersVector"`
	HomingMax           int       `json:"homingDaggersMax"`
	HomingMaxTime       float64   `json:"homingDaggersMaxTime"`
	DaggersFired        int       `json:"daggersFired"`
	DaggersFiredSlice   []int     `json:"daggersFiredVector"`
	DaggersHit          int       `json:"daggersHit"`
	DaggersHitSlice     []int     `json:"daggersHitVector"`
	EnemiesAlive        int       `json:"enemiesAlive"`
	EnemiesAliveSlice   []int     `json:"enemiesAliveVector"`
	EnemiesAliveMax     int       `json:"enemiesAliveMax"`
	EnemiesAliveMaxTime float64   `json:"enemiesAliveMaxTime"`
	EnemiesKilled       int       `json:"enemiesKilled"`
	EnemiesKilledSlice  []int     `json:"enemiesKilledVector"`
	DeathType           int       `json:"deathType"`
	ReplayPlayerID      int       `json:"replayPlayerID"`
	Version             string    `json:"version"`
	SurvivalHash        string    `json:"survivalHash"`
}
