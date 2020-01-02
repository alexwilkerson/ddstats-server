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
			state.Accuracy = roundToNearest(float64(state.DaggersHit)/float64(state.DaggersFired), 2)
		}
		states = append(states, &state)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return states, nil
}

// GetGems returns how many Gems in the game
func (g *GameModel) GetGems(id int) (int, error) {
	return 0, nil
}

// GetHomingDaggers returns how many homing daggers
func (g *GameModel) GetHomingDaggers(id int) (int, error) {
	return 0, nil
}

// GetAccuracy returns the game total accuracy
func (g *GameModel) GetAccuracy(id int) (int, error) {
	return 0, nil
}

// GetEnemiesAlive returns how many enemies are still alive
func (g *GameModel) GetEnemiesAlive(id int) (int, error) {
	return 0, nil
}

// GetEnemiesKilled return how many enemies had been killed
func (g *GameModel) GetEnemiesKilled(id int) (int, error) {
	return 0, nil
}

func roundToNearest(f float64, numberOfDecimalPlaces int) float64 {
	multiplier := math.Pow10(numberOfDecimalPlaces)
	return math.Round(f*multiplier) / multiplier
}
