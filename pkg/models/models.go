package models

import (
	"errors"
	"time"

	"gopkg.in/guregu/null.v3"
)

//ErrNoRecord will be returned when DB record not found
var ErrNoRecord = errors.New("no record found")
var ErrNoDiscordUserFound = errors.New("no entry associated with that discord ID")
var ErrDiscordUserVerified = errors.New("discord user is verified so cannot update their values")

//Game record representation
type Game struct {
	ID                   int         `json:"id" db:"id"`
	PlayerID             int         `json:"player_id" db:"player_id"`
	Granularity          int         `json:"granularity" db:"granularity"`
	GameTime             float64     `json:"game_time" db:"game_time"`
	DeathType            string      `json:"death_type" db:"death_type"`
	Gems                 int         `json:"gems" db:"gems"`
	HomingDaggers        int         `json:"homing_daggers" db:"homing_daggers"`
	DaggersFired         int         `json:"daggers_fired" db:"daggers_fired"`
	DaggersHit           int         `json:"daggers_hit" db:"daggers_hit"`
	Accuracy             float64     `json:"accuracy" db:"accuracy"`
	EnemiesAlive         int         `json:"enemies_alive" db:"enemies_alive"`
	EnemiesKilled        int         `json:"enemies_killed" db:"enemies_killed"`
	TimeStamp            time.Time   `json:"time_stamp" db:"time_stamp"`
	ReplayPlayerID       int         `json:"replay_player_id" db:"replay_player_id"`
	SurvivalHash         string      `json:"survival_hash" db:"survival_hash"`
	Version              null.String `json:"version" db:"version"`
	LevelTwoTime         float64     `json:"level_two_time" db:"level_two_time"`
	LevelThreeTime       float64     `json:"level_three_time" db:"level_three_time"`
	LevelFourTime        float64     `json:"level_four_time" db:"level_four_time"`
	HomingDaggersMaxTime float64     `json:"homing_daggers_max_time" db:"homing_daggers_max_time"`
	EnemiesAliveMaxTime  float64     `json:"enemies_alive_max_time" db:"enemies_alive_max_time"`
	HomingDaggersMax     int         `json:"homing_daggers_max" db:"homing_daggers_max"`
	EnemiesAliveMax      int         `json:"enemies_alive_max" db:"enemies_alive_max"`
}

// GameWithName is game with player_name included
type GameWithName struct {
	PlayerName string `json:"player_name" db:"player_name"`
	Game
}

// Player struct is for players
type Player struct {
	ID                     int     `json:"player_id" db:"id"`
	PlayerName             string  `json:"player_name" db:"player_name"`
	Rank                   int     `json:"rank" db:"rank"`
	GameTime               float64 `json:"game_time" db:"game_time"`
	DeathType              string  `json:"death_type" db:"death_type"`
	Gems                   int     `json:"gems" db:"gems"`
	DaggersHit             int     `json:"daggers_hit" db:"daggers_hit"`
	DaggersFired           int     `json:"daggers_fired" db:"daggers_fired"`
	EnemiesKilled          int     `json:"enemies_killed" db:"enemies_killed"`
	Accuracy               float64 `json:"accuracy" db:"accuracy"`
	OverallGameTime        float64 `json:"overall_game_time" db:"overall_game_time"`
	OverallAverageGameTime float64 `json:"overall_average_game_time" db:"overall_average_game_time"`
	OverallDeaths          int     `json:"overall_deaths" db:"overall_deaths"`
	OverallGems            int     `json:"overall_gems" db:"overall_gems"`
	OverallEnemiesKilled   int     `json:"overall_enemies_killed" db:"overall_enemies_killed"`
	OverallDaggersHit      int     `json:"overall_daggers_hit" db:"overall_daggers_hit"`
	OverallDaggersFired    int     `json:"overall_daggers_fired" db:"overall_daggers_fired"`
	OverallAccuracy        float64 `json:"overall_accuracy" db:"overall_accuracy"`
}

// State struct is for State
type State struct {
	GameID        int     `json:"game_id,omitempty" db:"game_id"`
	GameTime      float64 `json:"game_time" db:"game_time"`
	Gems          int     `json:"gems" db:"gems"`
	HomingDaggers int     `json:"homing_daggers" db:"homing_daggers"`
	DaggersHit    int     `json:"daggers_hit" db:"daggers_hit"`
	DaggersFired  int     `json:"daggers_fired" db:"daggers_fired"`
	Accuracy      float64 `json:"accuracy" db:"accuracy"`
	EnemiesAlive  int     `json:"enemies_alive" db:"enemies_alive"`
	EnemiesKilled int     `json:"enemies_killed" db:"enemies_killed"`
}

// Gems holds game time and gems
type Gems struct {
	GameTime float64 `json:"game_time" db:"game_time"`
	Gems     int     `json:"gems" db:"gems"`
}

// HomingDaggers holds game time and homing daggers
type HomingDaggers struct {
	GameTime      float64 `json:"game_time" db:"game_time"`
	HomingDaggers int     `json:"homing_daggers" db:"homing_daggers"`
}

// DaggersHit holds game time and daggers hit
type DaggersHit struct {
	GameTime   float64 `json:"game_time" db:"game_time"`
	DaggersHit int     `json:"daggers_hit" db:"daggers_hit"`
}

// DaggersFired holds game time and daggers fired
type DaggersFired struct {
	GameTime     float64 `json:"game_time" db:"game_time"`
	DaggersFired int     `json:"daggers_fired" db:"daggers_fired"`
}

// Accuracy holds game time and accuracy
type Accuracy struct {
	GameTime float64 `json:"game_time" db:"game_time"`
	Accuracy float64 `json:"accuracy" db:"accuracy"`
}

// EnemiesAlive holds game time and enemies alive
type EnemiesAlive struct {
	GameTime     float64 `json:"game_time" db:"game_time"`
	EnemiesAlive int     `json:"enemies_alive" db:"enemies_alive"`
}

// EnemiesKilled holds game time and enemies killed
type EnemiesKilled struct {
	GameTime      float64 `json:"game_time" db:"game_time"`
	EnemiesKilled int     `json:"enemies_killed" db:"enemies_killed"`
}

// SubmittedGame is used to decode the JSON struct that comes in when a player
// completes a game and is submitted
type SubmittedGame struct {
	PlayerID            int       `json:"playerID"`
	PlayerName          string    `json:"playerName"`
	Granularity         int       `json:"granularity"`
	GameTime            float64   `json:"inGameTimer"`
	GameTimeSlice       []float64 `json:"inGameTimerVector"`
	Gems                int       `json:"gems"`
	GemsSlice           []int     `json:"gemsVector"`
	LevelTwoTime        float64   `json:"levelTwoTime"`
	LevelThreeTime      float64   `json:"levelThreeTime"`
	LevelFourTime       float64   `json:"levelFourTime"`
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

// MOTD is Message of the Day and it's used by the client
// to display the message of the day
type MOTD struct {
	ID      int       `json:"id" db:"id"`
	Created time.Time `json:"created" db:"created"`
	Message string    `json:"motd" db:"message"`
}

type DiscordUser struct {
	DiscordID string `db:"discord_id"`
	DDID      int    `db:"dd_id"`
	Verified  bool   `db:"verified"`
}
