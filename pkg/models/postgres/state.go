package postgres

import (
	"database/sql"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

type StateModel struct {
	DB *sql.DB
}

func (s *StateModel) Insert(state *models.State) error {
	stmt := `INSERT INTO state(
			game_id,
			game_time,
			gems,
			homing_daggers,
			daggers_hit,
			daggers_fired,
			enemies_alive,
			enemies_killed
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
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
