package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// StateModel wraps the database connection
type StateModel struct {
	DB *sqlx.DB
}

// Insert inserts a state into the state table
func (s *StateModel) Insert(state *models.State) error {
	stmt := `
		INSERT INTO state(
			game_id,
			game_time,
			gems,
			homing_daggers,
			daggers_hit,
			daggers_fired,
			enemies_alive,
			enemies_killed) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.DB.Exec(stmt,
		state.GameID,
		state.GameTime,
		state.Gems,
		state.HomingDaggers,
		state.DaggersHit,
		state.DaggersFired,
		state.EnemiesAlive,
		state.EnemiesKilled,
	)
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a state into the state table
func (s *StateModel) InsertGRPC(state *models.State) error {
	stmt := `
		INSERT INTO state(
			game_id,
			game_time,
			gems,
			homing_daggers,
			daggers_hit,
			daggers_fired,
			enemies_alive,
			enemies_killed,
			total_gems,
			level_gems,
			gems_despawned,
			gems_eaten,
			daggers_eaten,
			per_enemy_alive_count,
			per_enemy_kill_count) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
				$11, $12, $13, $14, $15)`
	_, err := s.DB.Exec(stmt,
		state.GameID,
		state.GameTime,
		state.Gems,
		state.HomingDaggers,
		state.DaggersHit,
		state.DaggersFired,
		state.EnemiesAlive,
		state.EnemiesKilled,
		state.TotalGems,
		state.LevelGems,
		state.GemsDespawned,
		state.DaggersEaten,
		state.DaggersFired,
		pq.Array(state.PerEnemyAliveCount),
		pq.Array(state.PerEnemyKillCount),
	)
	if err != nil {
		return err
	}
	return nil
}

// GetAll returns a slice of states including all of the data from each state
func (s *StateModel) GetAll(id int) ([]*models.State, error) {
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
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return states, nil
}

// GetGems returns a slice game time and gems from the given game
func (s *StateModel) GetGems(id int) ([]*models.Gems, error) {
	var states []*models.Gems
	stmt := `
		SELECT round(game_time, 4) as game_time, gems
		FROM state
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetHomingDaggers returns a slice game time and homing daggers from the given game
func (s *StateModel) GetHomingDaggers(id int) ([]*models.HomingDaggers, error) {
	var states []*models.HomingDaggers
	stmt := `
		SELECT round(game_time, 4) as game_time, homing_daggers
		FROM state
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetDaggersHit returns a slice game time and daggers hit from the given game
func (s *StateModel) GetDaggersHit(id int) ([]*models.DaggersHit, error) {
	var states []*models.DaggersHit
	stmt := `
		SELECT round(game_time, 4) as game_time, daggers_hit
		FROM state
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetDaggersFired returns a slice game time and daggers fired from the given game
func (s *StateModel) GetDaggersFired(id int) ([]*models.DaggersFired, error) {
	var states []*models.DaggersFired
	stmt := `
		SELECT round(game_time, 4) as game_time, daggers_fired
		FROM state
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetAccuracy returns a slice game time and accuracy from the given game
func (s *StateModel) GetAccuracy(id int) ([]*models.Accuracy, error) {
	var states []*models.Accuracy
	stmt := `
		SELECT round(game_time, 4) as game_time, round(divzero(daggers_hit, daggers_fired)*100, 2) as accuracy
		FROM state
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetEnemiesAlive returns a slice game time and enemies alive from the given game
func (s *StateModel) GetEnemiesAlive(id int) ([]*models.EnemiesAlive, error) {
	var states []*models.EnemiesAlive
	stmt := `
		SELECT round(game_time, 4) as game_time, enemies_alive
		FROM state
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}

// GetEnemiesKilled returns a slice game time and enemies killed from the given game
func (s *StateModel) GetEnemiesKilled(id int) ([]*models.EnemiesKilled, error) {
	var states []*models.EnemiesKilled
	stmt := `
		SELECT round(game_time, 4) as game_time, enemies_killed
		FROM state
		WHERE game_id=$1
		ORDER BY game_time ASC`
	err := s.DB.Select(&states, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return states, nil
}
