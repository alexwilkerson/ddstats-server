package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type CollectorPlayerModel struct {
	DB *sqlx.DB
}

func (cpm *CollectorPlayerModel) Select(playerID int) (*models.CollectorPlayer, error) {
	var collectorPlayer models.CollectorPlayer
	stmt := `
		SELECT *
		FROM collector_player
		WHERE id=$1`
	err := cpm.DB.Get(&collectorPlayer, stmt, playerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return &collectorPlayer, nil
}

func (cpm *CollectorPlayerModel) NewPlayer(tx *sqlx.Tx, collectorPlayerID int) error {
	stmt := `
		INSERT INTO collector_player(id)
		VALUES ($1)
		ON CONFLICT (id)
		DO NOTHING`
	_, err := tx.Exec(stmt, collectorPlayerID)
	if err != nil {
		return err
	}
	return nil
}

func (cpm *CollectorPlayerModel) UpsertPlayer(tx *sqlx.Tx, player *ddapi.Player, collectorRunID int, lastActive *time.Time) error {
	stmt := `
		INSERT INTO collector_player(
			id,
			player_name,
			last_active,
			rank,
			game_time,
			death_type,
			gems,
			daggers_hit,
			daggers_fired,
			enemies_killed,
			overall_game_time,
			overall_deaths,
			overall_gems,
			overall_enemies_killed,
			overall_daggers_hit,
			overall_daggers_fired
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (id) DO
		UPDATE SET
			id=$1,
			player_name=$2,
			last_active=$3,
			rank=$4,
			game_time=$5,
			death_type=$6,
			gems=$7,
			daggers_hit=$8,
			daggers_fired=$9,
			enemies_killed=$10,
			overall_game_time=$11,
			overall_deaths=$12,
			overall_gems=$13,
			overall_enemies_killed=$14,
			overall_daggers_hit=$15,
			overall_daggers_fired=$16
		WHERE collector_player.id=$1`
	_, err := tx.Exec(stmt,
		player.PlayerID,
		player.PlayerName,
		lastActive,
		player.Rank,
		player.GameTime,
		player.DeathType,
		player.Gems,
		player.DaggersHit,
		player.DaggersFired,
		player.EnemiesKilled,
		player.OverallGameTime,
		player.OverallDeaths,
		player.OverallGems,
		player.OverallEnemiesKilled,
		player.OverallDaggersHit,
		player.OverallDaggersFired,
	)
	if err != nil {
		fmt.Println("upsert:", err)
		return err
	}
	return nil
}
