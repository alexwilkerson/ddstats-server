package postgres

import (
	"database/sql"
	"errors"
	"math"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

// GameModel wraps database connection
type GameModel struct {
	DB *sql.DB
}

// Get retreives the entire game obeject
func (g *GameModel) Get(id int) (*models.Game, error) {
	var gameModel models.Game

	stmt := `SELECT * FROM game WHERE id=$1`
	err := g.DB.QueryRow(stmt, id).Scan(
		&gameModel.ID,
		&gameModel.PlayerID,
		&gameModel.Granularity,
		&gameModel.GameTime,
		&gameModel.DeathType,
		&gameModel.Gems,
		&gameModel.HomingDaggers,
		&gameModel.DaggersFired,
		&gameModel.DaggersHit,
		&gameModel.EnemiesAlive,
		&gameModel.EnemiesKilled,
		&gameModel.TimeStamp,
		&gameModel.ReplayPlayerID,
		&gameModel.SurvivalHash,
		&gameModel.Version,
		&gameModel.LevelTwoTime,
		&gameModel.LevelThreeTime,
		&gameModel.LevelFourTime,
		&gameModel.HomingDaggersMaxTime,
		&gameModel.EnemiesAliveMaxTime,
		&gameModel.HomingDaggersMax,
		&gameModel.EnemiesAliveMax,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return &gameModel, nil
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

func roundToNearest(f float64, numberOfDecimalPlaces int) float64 {
	multiplier := math.Pow10(numberOfDecimalPlaces)
	return math.Round(f*multiplier) / multiplier
}
