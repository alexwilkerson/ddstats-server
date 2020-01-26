package postgres

import (
	"github.com/jmoiron/sqlx"
)

type ReplayPlayerModel struct {
	DB *sqlx.DB
}

func (p *ReplayPlayerModel) Upsert(playerID int, playerName string) error {
	stmt := `
		INSERT INTO replay_player(
			id,
			player_name
		) VALUES ($1, $2)
		ON CONFLICT (id) DO
		UPDATE SET
			player_name=$2
		WHERE replay_player.id=$1`
	_, err := p.DB.Exec(stmt,
		playerID,
		playerName,
	)
	if err != nil {
		return err
	}
	return nil
}
