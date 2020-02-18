package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type CollectorHighScoreModel struct {
	DB *sqlx.DB
}

func (crsm *CollectorHighScoreModel) Insert(tx *sqlx.Tx, collectorRunID, collectorPlayerID int, score float64) error {
	stmt := `
		INSERT INTO collector_high_score(collector_run_id, collector_player_id, score)
		VALUES ($1, $2, $3)`
	_, err := tx.Exec(stmt, collectorRunID, collectorPlayerID, score)
	if err != nil {
		return err
	}
	return nil
}

func (crsm *CollectorHighScoreModel) Select(collectorRunID int) ([]*models.CollectorHighScore, error) {
	scores := []*models.CollectorHighScore{}
	stmt := `
		SELECT
			collector_run_id,
			collector_player_id,
			collector_player.player_name AS collector_player_name,
			collector_player.rank AS collector_player_rank,
			score
		FROM collector_high_score
		JOIN collector_player ON collector_player_id=collector_player.id
		WHERE collector_run_id=$1
		ORDER BY score DESC`
	err := crsm.DB.Select(&scores, stmt, collectorRunID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return scores, nil
}
