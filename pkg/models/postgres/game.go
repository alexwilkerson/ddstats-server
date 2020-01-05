package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
	"github.com/jmoiron/sqlx"
)

// GameModel wraps database connection
type GameModel struct {
	DB *sqlx.DB
}

const (
	v3SurvivalHashA = "5ff43e37d0f85e068caab5457305754e"
	v3SurvivalHashB = "569fead87abf4d30fdee4231a6398051"
)

// GetTop retreives a slice of the top games in the database with a given limit
func (g *GameModel) GetTop(limit int) ([]*models.GameWithName, error) {
	var games []*models.GameWithName

	stmt := fmt.Sprintf(`
		SELECT
			game.id,
			player_id,
			player_name,
			granularity,
			round(game.game_time, 4) as game_time,
			death_type.name as death_type,
			game.gems,
			game.homing_daggers,
			game.daggers_fired,
			game.daggers_hit,
			round(divzero(game.daggers_hit, game.daggers_fired)*100, 2) as accuracy,
			game.enemies_alive,
			game.enemies_killed,
			time_stamp,
			replay_player_id,
			survival_hash,
			version,
			level_two_time,
			level_three_time,
			level_four_time,
			homing_daggers_max_time,
			enemies_alive_max_time,
			homing_daggers_max,
			enemies_alive_max
		FROM game JOIN player ON game.player_id=player.id JOIN death_type ON game.death_type=death_type.id
		WHERE replay_player_id=0 AND (survival_hash='%s' OR survival_hash='%s')
		ORDER BY game_time DESC LIMIT %d`, v3SurvivalHashA, v3SurvivalHashB, limit)
	err := g.DB.Select(&games, stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	sort.Sort(byTime(games))

	return games, nil
}

type byTime []*models.GameWithName

func (a byTime) Len() int           { return len(a) }
func (a byTime) Less(i, j int) bool { return (*a[i]).GameTime < (*a[j]).GameTime }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// GetRecent retreives a slice of users using a specified page size and page num starting at 1
func (g *GameModel) GetRecent(playerID, pageSize, pageNum int) ([]*models.GameWithName, error) {
	var where string
	if playerID != 0 {
		where = fmt.Sprintf("WHERE game.player_id=$1 AND game.replay_player_id=0")
	}

	var games []*models.GameWithName

	stmt := fmt.Sprintf(`
		SELECT
			game.id,
			player_id,
			player_name,
			granularity,
			round(game.game_time, 4) as game_time,
			death_type.name as death_type,
			game.gems,
			game.homing_daggers,
			game.daggers_fired,
			game.daggers_hit,
			round(divzero(game.daggers_hit, game.daggers_fired)*100, 2) as accuracy,
			game.enemies_alive,
			game.enemies_killed,
			time_stamp,
			replay_player_id,
			survival_hash,
			version,
			level_two_time,
			level_three_time,
			level_four_time,
			homing_daggers_max_time,
			enemies_alive_max_time,
			homing_daggers_max,
			enemies_alive_max
		FROM game JOIN player ON game.player_id=player.id JOIN death_type ON game.death_type=death_type.id %s
		ORDER BY id DESC LIMIT %d OFFSET %d`, where, pageSize, (pageNum-1)*pageSize)
	var err error
	if playerID != 0 {
		err = g.DB.Select(&games, stmt, playerID)
	} else {
		err = g.DB.Select(&games, stmt)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return games, nil
}

// Get retreives the entire game obeject
func (g *GameModel) Get(id int) (*models.GameWithName, error) {
	var game models.GameWithName
	stmt := `
		SELECT
			game.id,
			player_id,
			player_name,
			granularity,
			round(game.game_time, 4) as game_time,
			death_type.name as death_type,
			game.gems,
			game.homing_daggers,
			game.daggers_fired,
			game.daggers_hit,
			round(divzero(game.daggers_hit, game.daggers_fired)*100, 2) as accuracy,
			game.enemies_alive,
			game.enemies_killed,
			time_stamp,
			replay_player_id,
			survival_hash,
			version,
			level_two_time,
			level_three_time,
			level_four_time,
			homing_daggers_max_time,
			enemies_alive_max_time,
			homing_daggers_max,
			enemies_alive_max
		FROM game JOIN player ON game.player_id=player.id JOIN death_type ON game.death_type=death_type.id
		WHERE game.id=$1`
	err := g.DB.Get(&game, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return &game, nil
}

// GetAll returns a slice of states including all of the data from each state
func (g *GameModel) GetAll(id int) ([]*models.State, error) {
	var states []*models.State
	stmt := `
		SELECT
			round(game_time, 4) as game_time,
			gems,
			homing_daggers,
			daggers_hit,
			daggers_fired,
			round(divzero(daggers_hit, daggers_fired)*100, 2) as accuracy,
			enemies_alive,
			enemies_killed
		FROM state
		WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return states, nil
}

// GetGems returns a slice game time and gems from the given game
func (g *GameModel) GetGems(id int) ([]*models.Gems, error) {
	var states []*models.Gems
	stmt := `SELECT round(game_time, 4) as game_time, gems
			 FROM state
			 WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetHomingDaggers returns a slice game time and homing daggers from the given game
func (g *GameModel) GetHomingDaggers(id int) ([]*models.HomingDaggers, error) {
	var states []*models.HomingDaggers
	stmt := `SELECT round(game_time, 4) as game_time, homing_daggers
			 FROM state
			 WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetDaggersHit returns a slice game time and daggers hit from the given game
func (g *GameModel) GetDaggersHit(id int) ([]*models.DaggersHit, error) {
	var states []*models.DaggersHit
	stmt := `SELECT round(game_time, 4) as game_time, daggers_hit
			 FROM state
			 WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetDaggersFired returns a slice game time and daggers fired from the given game
func (g *GameModel) GetDaggersFired(id int) ([]*models.DaggersFired, error) {
	var states []*models.DaggersFired
	stmt := `SELECT round(game_time, 4) as game_time, daggers_fired
			 FROM state
			 WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetAccuracy returns a slice game time and accuracy from the given game
func (g *GameModel) GetAccuracy(id int) ([]*models.Accuracy, error) {
	var states []*models.Accuracy
	stmt := `SELECT round(game_time, 4) as game_time, round(divzero(daggers_hit, daggers_fired)*100, 2) as accuracy
			 FROM state
			 WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetEnemiesAlive returns a slice game time and enemies alive from the given game
func (g *GameModel) GetEnemiesAlive(id int) ([]*models.EnemiesAlive, error) {
	var states []*models.EnemiesAlive
	stmt := `SELECT round(game_time, 4) as game_time, enemies_alive
			 FROM state
			 WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetEnemiesKilled returns a slice game time and enemies killed from the given game
func (g *GameModel) GetEnemiesKilled(id int) ([]*models.EnemiesKilled, error) {
	var states []*models.EnemiesKilled
	stmt := `SELECT round(game_time, 4) as game_time, enemies_killed
			 FROM state
			 WHERE game_id=$1`
	err := g.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetTotalCount returns the total number of games in the database
func (g *GameModel) GetTotalCount(playerID int) (int, error) {
	var err error
	var stmt string
	var gameCount int
	if playerID != 0 {
		stmt = "SELECT COUNT(1) FROM game WHERE player_id=$1 AND replay_player_id=0"
		err = g.DB.QueryRow(stmt, playerID).Scan(&gameCount)
	} else {
		stmt = "SELECT COUNT(1) FROM game"
		err = g.DB.QueryRow(stmt).Scan(&gameCount)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}
	return gameCount, nil
}

func roundToNearest(f float64, numberOfDecimalPlaces int) float64 {
	multiplier := math.Pow10(numberOfDecimalPlaces)
	return math.Round(f*multiplier) / multiplier
}
