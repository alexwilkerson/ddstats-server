package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type CollectorActivePlayerModel struct {
	DB *sqlx.DB
}

func (cap *CollectorActivePlayerModel) Insert(tx *sqlx.Tx, runID, playerID, rank, rankImprovement int, gameTime, gameTimeImprovement float64) error {
	stmt := `
		INSERT INTO collector_active_player (collector_run_id, collector_player_id, rank, rank_improvement, game_time, game_time_improvement)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := tx.Exec(stmt, runID, playerID, rank, rankImprovement, gameTime, gameTimeImprovement)
	if err != nil {
		return err
	}
	return nil
}

func (cap *CollectorActivePlayerModel) Select(runID int) ([]*models.CollectorActivePlayer, error) {
	players := []*models.CollectorActivePlayer{}
	stmt := `
		SELECT
			collector_run_id,
			collector_player_id,
			player.player_name AS collector_player_name,
			collector_active_player.rank,
			rank_improvement,
			collector_active_player.game_time,
			collector_active_player.game_time_improvement
		FROM collector_active_player JOIN player ON collector_player_id=player.id
		WHERE collector_run_id=$1`
	err := cap.DB.Select(&players, stmt, runID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return players, nil
}
