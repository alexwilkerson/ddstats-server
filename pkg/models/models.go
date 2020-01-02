package models

import (
	"database/sql"
	"encoding/json"
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
	Version              JsonNullString `json:"version"`
	LevelTwoTime         float64        `json:"level_two_time"`
	LevelThreeTime       float64        `json:"level_three_time"`
	LevelFourTime        float64        `json:"level_four_time"`
	HomingDaggersMaxTime float64        `json:"homing_daggers_max_time"`
	EnemiesAliveMaxTime  float64        `json:"enemies_alive_max_time"`
	HomingDaggersMax     uint           `json:"homing_daggers_max"`
	EnemiesAliveMax      uint           `json:"enemies_alive_max"`
}

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

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}
