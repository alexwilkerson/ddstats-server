package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type CollectorNewPlayerModel struct {
	DB *sqlx.DB
}

func (cnp *CollectorNewPlayerModel) Insert(tx *sqlx.Tx, runID, playerID, rank int, gameTime float64) error {
	stmt := `
		INSERT INTO collector_new_player (collector_run_id, collector_player_id, rank, game_time)
		VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(stmt, runID, playerID, rank, gameTime)
	if err != nil {
		return err
	}
	return nil
}

func (cnp *CollectorNewPlayerModel) Select(runID int) ([]*models.CollectorNewPlayer, error) {
	players := []*models.CollectorNewPlayer{}
	stmt := `
		SELECT collector_run_id, collector_player_id, player.player_name AS collector_player_name, game_time
		FROM collector_new_player JOIN player ON collector_new_id=player.id
		WHERE collector_run_id=$1`
	err := cnp.DB.Select(&players, stmt, runID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return players, nil
}
