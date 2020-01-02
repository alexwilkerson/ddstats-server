package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

// PlayerModel wraps the database connection
type PlayerModel struct {
	DB *sql.DB
}

// Get returns a single player record
func (p *PlayerModel) Get(id int) (*models.Player, error) {
	var player models.Player

	stmt := "SELECT * FROM player WHERE id=$1"
	err := p.DB.QueryRow(stmt, id).Scan(
		&player.ID,
		&player.PlayerName,
		&player.Rank,
		&player.GameTime,
		&player.DeathType,
		&player.Gems,
		&player.DaggersHit,
		&player.DaggersFired,
		&player.EnemiesKilled,
		&player.Accuracy,
		&player.OverallTime,
		&player.OverallDeaths,
		&player.OverallGems,
		&player.OverallEnemiesKilled,
		&player.OverallDaggersHit,
		&player.OverallDaggersFired,
		&player.OverallAccuracy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return &player, nil
}

func (p *PlayerModel) GetAll(pageSize, pageNum int) ([]*models.Player, error) {
	var players []*models.Player

	stmt := fmt.Sprintf("SELECT * FROM player WHERE id<>-1 ORDER BY game_time DESC LIMIT %d OFFSET %d", pageSize, (pageNum-1)*pageSize)
	rows, err := p.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var player models.Player
		err = rows.Scan(
			&player.ID,
			&player.PlayerName,
			&player.Rank,
			&player.GameTime,
			&player.DeathType,
			&player.Gems,
			&player.DaggersHit,
			&player.DaggersFired,
			&player.EnemiesKilled,
			&player.Accuracy,
			&player.OverallTime,
			&player.OverallDeaths,
			&player.OverallGems,
			&player.OverallEnemiesKilled,
			&player.OverallDaggersHit,
			&player.OverallDaggersFired,
			&player.OverallAccuracy,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, &player)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return players, nil
}
