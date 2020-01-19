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

func (cap *CollectorActivePlayerModel) Insert(tx *sqlx.Tx, runID, playerID, rank, rankImprovement int, gameTime, gameTimeImprovement, sinceGameTime float64, sinceDeaths int) error {
	stmt := `
		INSERT INTO collector_active_player (collector_run_id, collector_player_id, rank, rank_improvement, game_time, game_time_improvement, since_game_time, since_deaths)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := tx.Exec(stmt, runID, playerID, rank, rankImprovement, gameTime, gameTimeImprovement, sinceGameTime, sinceDeaths)
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
			collector_player.player_name AS collector_player_name,
			collector_active_player.rank,
			rank_improvement,
			collector_active_player.game_time,
			ROUND(collector_active_player.game_time_improvement, 4) AS game_time_improvement,
			ROUND(since_game_time, 4) AS since_game_time,
			since_deaths
		FROM collector_active_player JOIN collector_player ON collector_player_id=collector_player.id
		WHERE collector_run_id=$1
		ORDER BY collector_active_player.rank ASC`
	err := cap.DB.Select(&players, stmt, runID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return players, nil
}
