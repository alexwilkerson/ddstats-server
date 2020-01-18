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

func (cap *CollectorActivePlayerModel) Insert(tx *sqlx.Tx, runID, playerID, rank int, gameTime float64) error {
	stmt := `
		INSERT INTO collector_active_player (collector_run_id, collector_player_id, rank, game_time)
		VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(stmt, runID, playerID, rank, gameTime)
	if err != nil {
		return err
	}
	return nil
}

func (cap *CollectorActivePlayerModel) Select(runID int) ([]*models.CollectorActivePlayer, error) {
	players := []*models.CollectorActivePlayer{}
	stmt := `
		SELECT collector_run_id, collector_player_id, player.player_name AS collector_player_name, game_time
		FROM collector_active_player JOIN player ON collector_player_id=player.id
		WHERE collector_run_id=$1`
	err := cap.DB.Select(&players, stmt, runID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return players, nil
}
