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

// GetAll retreives a slice of users using a specified page size and page num starting at 1
func (g *GameModel) GetTop(limit int) ([]*models.Game, error) {
	var games []*models.Game

	stmt := fmt.Sprintf(`SELECT *
						 FROM game 
						 WHERE survival_hash='%s' OR survival_hash='%s'
						 ORDER BY game_time DESC LIMIT %d`, v3SurvivalHashA, v3SurvivalHashB, limit)
	rows, err := g.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var game models.Game
		err = rows.Scan(
			&game.ID,
			&game.PlayerID,
			&game.Granularity,
			&game.GameTime,
			&game.DeathType,
			&game.Gems,
			&game.HomingDaggers,
			&game.DaggersFired,
			&game.DaggersHit,
			&game.EnemiesAlive,
			&game.EnemiesKilled,
			&game.TimeStamp,
			&game.ReplayPlayerID,
			&game.SurvivalHash,
			&game.Version,
			&game.LevelTwoTime,
			&game.LevelThreeTime,
			&game.LevelFourTime,
			&game.HomingDaggersMaxTime,
			&game.EnemiesAliveMaxTime,
			&game.HomingDaggersMax,
			&game.EnemiesAliveMax,
		)
		if err != nil {
			return nil, err
		}
		if game.DaggersFired > 0 {
			game.Accuracy = roundToNearest(float64(game.DaggersHit)/float64(game.DaggersFired)*100, 2)
		}
		games = append(games, &game)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	sort.Sort(byTime(games))

	return games, nil
}

type byTime []*models.Game

func (a byTime) Len() int           { return len(a) }
func (a byTime) Less(i, j int) bool { return (*a[i]).GameTime < (*a[j]).GameTime }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// GetRecent retreives a slice of users using a specified page size and page num starting at 1
func (g *GameModel) GetRecent(pageSize, pageNum int) ([]*models.Game, error) {
	var games []*models.Game

	stmt := fmt.Sprintf("SELECT * FROM game ORDER BY time_stamp DESC LIMIT %d OFFSET %d", pageSize, (pageNum-1)*pageSize)
	rows, err := g.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var game models.Game
		err = rows.Scan(
			&game.ID,
			&game.PlayerID,
			&game.Granularity,
			&game.GameTime,
			&game.DeathType,
			&game.Gems,
			&game.HomingDaggers,
			&game.DaggersFired,
			&game.DaggersHit,
			&game.EnemiesAlive,
			&game.EnemiesKilled,
			&game.TimeStamp,
			&game.ReplayPlayerID,
			&game.SurvivalHash,
			&game.Version,
			&game.LevelTwoTime,
			&game.LevelThreeTime,
			&game.LevelFourTime,
			&game.HomingDaggersMaxTime,
			&game.EnemiesAliveMaxTime,
			&game.HomingDaggersMax,
			&game.EnemiesAliveMax,
		)
		if err != nil {
			return nil, err
		}
		if game.DaggersFired > 0 {
			game.Accuracy = roundToNearest(float64(game.DaggersHit)/float64(game.DaggersFired)*100, 2)
		}
		games = append(games, &game)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return games, nil
}

// Get retreives the entire game obeject
func (g *GameModel) Get(id int) (*models.Game, error) {
	var game models.Game

	stmt := `SELECT * FROM game WHERE id=$1`
	err := g.DB.QueryRow(stmt, id).Scan(
		&game.ID,
		&game.PlayerID,
		&game.Granularity,
		&game.GameTime,
		&game.DeathType,
		&game.Gems,
		&game.HomingDaggers,
		&game.DaggersFired,
		&game.DaggersHit,
		&game.EnemiesAlive,
		&game.EnemiesKilled,
		&game.TimeStamp,
		&game.ReplayPlayerID,
		&game.SurvivalHash,
		&game.Version,
		&game.LevelTwoTime,
		&game.LevelThreeTime,
		&game.LevelFourTime,
		&game.HomingDaggersMaxTime,
		&game.EnemiesAliveMaxTime,
		&game.HomingDaggersMax,
		&game.EnemiesAliveMax,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	if game.DaggersFired > 0 {
		game.Accuracy = roundToNearest(float64(game.DaggersHit)/float64(game.DaggersFired)*100, 2)
	}
	return &game, nil
}

func (g *GameModel) GetAll(id int) ([]*models.State, error) {
	stmt := `SELECT game_time, gems, homing_daggers, daggers_hit, daggers_fired, enemies_alive, enemies_killed
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.State
	for rows.Next() {
		var state models.State
		err = rows.Scan(
			&state.GameTime,
			&state.Gems,
			&state.HomingDaggers,
			&state.DaggersHit,
			&state.DaggersFired,
			&state.EnemiesAlive,
			&state.EnemiesKilled,
		)
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		if state.DaggersFired > 0 {
			state.Accuracy = roundToNearest(float64(state.DaggersHit)/float64(state.DaggersFired)*100, 2)
		}
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (g *GameModel) GetGems(id int) ([]*models.Gems, error) {
	stmt := `SELECT game_time, gems
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.Gems
	for rows.Next() {
		var state models.Gems
		err = rows.Scan(
			&state.GameTime,
			&state.Gems,
		)
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (g *GameModel) GetHomingDaggers(id int) ([]*models.HomingDaggers, error) {
	stmt := `SELECT game_time, homing_daggers
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.HomingDaggers
	for rows.Next() {
		var state models.HomingDaggers
		err = rows.Scan(
			&state.GameTime,
			&state.HomingDaggers,
		)
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (g *GameModel) GetDaggersHit(id int) ([]*models.DaggersHit, error) {
	stmt := `SELECT game_time, daggers_hit
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.DaggersHit
	for rows.Next() {
		var state models.DaggersHit
		err = rows.Scan(
			&state.GameTime,
			&state.DaggersHit,
		)
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (g *GameModel) GetDaggersFired(id int) ([]*models.DaggersFired, error) {
	stmt := `SELECT game_time, daggers_fired
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.DaggersFired
	for rows.Next() {
		var state models.DaggersFired
		err = rows.Scan(
			&state.GameTime,
			&state.DaggersFired,
		)
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (g *GameModel) GetAccuracy(id int) ([]*models.Accuracy, error) {
	stmt := `SELECT game_time, daggers_hit, daggers_fired
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.Accuracy
	for rows.Next() {
		var state models.Accuracy
		var daggersHit, daggersFired int
		err = rows.Scan(
			&state.GameTime,
			&daggersHit,
			&daggersFired,
		)
		if daggersFired > 0 {
			state.Accuracy = roundToNearest(float64(daggersHit)/float64(daggersFired)*100, 2)
		}
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (g *GameModel) GetEnemiesAlive(id int) ([]*models.EnemiesAlive, error) {
	stmt := `SELECT game_time, enemies_alive
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.EnemiesAlive
	for rows.Next() {
		var state models.EnemiesAlive
		err = rows.Scan(
			&state.GameTime,
			&state.EnemiesAlive,
		)
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (g *GameModel) GetEnemiesKilled(id int) ([]*models.EnemiesKilled, error) {
	stmt := `SELECT game_time, enemies_killed
			 FROM state
			 WHERE game_id=$1`
	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var states []*models.EnemiesKilled
	for rows.Next() {
		var state models.EnemiesKilled
		err = rows.Scan(
			&state.GameTime,
			&state.EnemiesKilled,
		)
		if err != nil {
			return nil, err
		}
		state.GameTime = roundToNearest(state.GameTime, 4)
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

// GetTotalCount returns the total number of games in the database
func (g *GameModel) GetTotalCount() (int, error) {
	var gameCount int
	stmt := "SELECT COUNT(1) FROM game"
	err := g.DB.QueryRow(stmt).Scan(&gameCount)
	if err != nil {
		return 0, err
	}
	return gameCount, nil
}

func roundToNearest(f float64, numberOfDecimalPlaces int) float64 {
	multiplier := math.Pow10(numberOfDecimalPlaces)
	return math.Round(f*multiplier) / multiplier
}
