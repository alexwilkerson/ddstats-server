package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
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
	ID        int       `json:"id" db:"id"`
	TimeStamp time.Time `json:"created" db:"time_stamp"`
	Message   string    `json:"motd" db:"message"`
}

type DiscordUser struct {
	DiscordID string `db:"discord_id"`
	DDID      int    `db:"dd_id"`
	Verified  bool   `db:"verified"`
}

type Release struct {
	Version   string    `json:"version" db:"version"`
	TimeStamp time.Time `json:"time_stamp" db:"time_stamp"`
	Body      string    `json:"body" db:"body"`
	FileName  string    `json:"file_name" db:"file_name"`
}

type News struct {
	ID        int       `json:"id" db:"id"`
	TimeStamp time.Time `json:"time_stamp" db:"time_stamp"`
	Title     string    `json:"title" db:"title"`
	Body      string    `json:"body" db:"body"`
}

type CollectorRun struct {
	ID                                  int       `json:"-" db:"id"`
	TimeStamp                           time.Time `json:"time_stamp" db:"time_stamp"`
	RunTime                             Duration  `json:"-" db:"run_time"`
	GlobalPlayers                       int       `json:"global_players" db:"global_players"`
	NewPlayers                          int       `json:"new_players" db:"new_players"`
	ActivePlayers                       int       `json:"active_players" db:"active_players"`
	InactivePlayers                     int       `json:"inactive_players" db:"inactive_players"`
	PlayersWithNewScores                int       `json:"players_with_new_scores" db:"players_with_new_scores"`
	PlayersWithNewRanks                 int       `json:"players_with_new_ranks" db:"players_with_new_ranks"`
	AverageImprovementTime              float64   `json:"average_improvement_time" db:"average_improvement_time"`
	AverageRankImprovement              float64   `json:"average_rank_improvement" db:"average_rank_improvement"`
	AverageGameTimePerActivePlayer      float64   `json:"average_game_time_per_active_player" db:"average_game_time_per_active_player"`
	AverageDeathsPerActivePlayer        float64   `json:"average_deaths_per_active_player" db:"average_deaths_per_active_player"`
	AverageGemsPerActivePlayer          float64   `json:"average_gems_per_active_player" db:"average_gems_per_active_player"`
	AverageEnemiesKilledPerActivePlayer float64   `json:"average_enemies_killed_per_active_player" db:"average_enemies_killed_per_active_player"`
	AverageDaggersHitPerActivePlayer    float64   `json:"average_daggers_hit_per_active_player" db:"average_daggers_hit_per_active_player"`
	AverageDaggersFiredPerActivePlayer  float64   `json:"average_daggers_fired_per_active_player" db:"average_daggers_fired_per_active_player"`
	AverageAccuracyPerActivePlayer      float64   `json:"average_accuracy_per_active_player" db:"average_accuracy_per_active_player"`
	GlobalGameTime                      float64   `json:"global_game_time" db:"global_game_time"`
	GlobalDeaths                        int       `json:"global_deaths" db:"global_deaths"`
	GlobalGems                          int       `json:"global_gems" db:"global_gems"`
	GlobalEnemiesKilled                 int       `json:"global_enemies_killed" db:"global_enemies_killed"`
	GlobalDaggersHit                    int       `json:"global_daggers_hit" db:"global_daggers_hit"`
	GlobalDaggersFired                  int       `json:"global_daggers_fired" db:"global_daggers_fired"`
	GlobalAccuracy                      float64   `json:"global_accuracy" db:"global_accuracy"`
	SinceGameTime                       float64   `json:"since_game_time" db:"since_game_time"`
	SinceDeaths                         int       `json:"since_deaths" db:"since_deaths"`
	SinceGems                           int       `json:"since_gems" db:"since_gems"`
	SinceEnemiesKilled                  int       `json:"since_enemies_killed" db:"since_enemies_killed"`
	SinceDaggersHit                     int       `json:"since_daggers_hit" db:"since_daggers_hit"`
	SinceDaggersFired                   int       `json:"since_daggers_fired" db:"since_daggers_fired"`
	SinceAccuracy                       float64   `json:"since_accuracy" db:"since_accuracy"`
	Fallen                              int       `json:"fallen" db:"fallen"`
	Swarmed                             int       `json:"swarmed" db:"swarmed"`
	Impaled                             int       `json:"impaled" db:"impaled"`
	Gored                               int       `json:"gored" db:"gored"`
	Infested                            int       `json:"infested" db:"infested"`
	Opened                              int       `json:"opened" db:"opened"`
	Purged                              int       `json:"purged" db:"purged"`
	Desecrated                          int       `json:"desecrated" db:"desecrated"`
	Sacrificed                          int       `json:"sacrificed" db:"sacrificed"`
	Eviscerated                         int       `json:"eviscerated" db:"eviscerated"`
	Annihilated                         int       `json:"annihilated" db:"annihilated"`
	Intoxicated                         int       `json:"intoxicated" db:"intoxicated"`
	Envenmonated                        int       `json:"envenmonated" db:"envenmonated"`
	Incarnated                          int       `json:"incarnated" db:"incarnated"`
	Discarnated                         int       `json:"discarnated" db:"discarnated"`
	Barbed                              int       `json:"barbed" db:"barbed"`
}

type CollectorPlayer struct {
	ID                   int     `db:"id"`
	PlayerName           string  `db:"player_name"`
	Rank                 int     `db:"rank"`
	GameTime             float64 `db:"game_time"`
	DeathType            string  `db:"death_type"`
	Gems                 int     `db:"gems"`
	DaggersHit           int     `db:"daggers_hit"`
	DaggersFired         int     `db:"daggers_fired"`
	EnemiesKilled        int     `db:"enemies_killed"`
	OverallGameTime      float64 `db:"overall_game_time"`
	OverallDeaths        int     `db:"overall_deaths"`
	OverallGems          int     `db:"overall_gems"`
	OverallEnemiesKilled int     `db:"overall_enemies_killed"`
	OverallDaggersHit    int     `db:"overall_daggers_hit"`
	OverallDaggersFired  int     `db:"overall_daggers_fired"`
}

type CollectorHighScore struct {
	CollectorRunID      int     `json:"-" db:"collector_run_id"`
	CollectorPlayerID   int     `json:"player_id" db:"collector_player_id"`
	CollectorPlayerName string  `json:"player_name" db:"collector_player_name"`
	Score               float64 `json:"score" db:"score"`
}

type CollectorActivePlayer struct {
	CollectorRunID      int     `json:"-" db:"collector_run_id"`
	CollectorPlayerID   int     `json:"player_id" db:"collector_player_id"`
	CollectorPlayerName string  `json:"player_name" db:"collector_player_name"`
	Rank                int     `json:"rank" db:"rank"`
	RankImprovement     int     `json:"rank_improvement" db:"rank_improvement"`
	GameTime            float64 `json:"game_time" db:"game_time"`
}

type CollectorNewPlayer struct {
	CollectorRunID      int     `json:"-" db:"collector_run_id"`
	CollectorPlayerID   int     `json:"player_id" db:"collector_player_id"`
	CollectorPlayerName string  `json:"player_name" db:"collector_player_name"`
	Rank                int     `json:"rank" db:"rank"`
	GameTime            float64 `json:"game_time" db:"game_time"`
}

type Duration time.Duration

func (d Duration) Value() (driver.Value, error) {
	return driver.Value(int64(d)), nil
}

func (d *Duration) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case int64:
		*d = Duration(v)
	case nil:
		*d = Duration(0)
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Duration from: %#v", v)
	}
	return nil
}
