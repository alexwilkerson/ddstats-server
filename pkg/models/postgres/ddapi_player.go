package postgres

import (
	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
	"github.com/jmoiron/sqlx"
)

// DDPlayersModel wraps the database connection
type DDPlayersModel struct {
	DB *sqlx.DB
}

// Insert takes the Player struct from the ddapi package and inserts it into
// the player table in the database
func (ddp *DDPlayersModel) Insert(player *ddapi.Player) error {
	stmt := `
		INSERT INTO player(
			player_name,
			rank,
			game_time,
			death_type,
			gems,
			daggers_hit,
			daggers_fired,
			enemies_killed,
			accuracy,
			overall_time,
			overall_deaths,
			overall_gems,
			overall_enemies_killed,
			overall_daggers_hit,
			overall_daggers_fired,
			overall_accuracy
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err := ddp.DB.Exec(stmt,
		player.PlayerName,
		player.Rank,
		player.GameTime,
		player.DeathType,
		player.Gems,
		player.DaggersHit,
		player.DaggersFired,
		player.EnemiesKilled,
		player.Accuracy,
		player.OverallTime,
		player.OverallDeaths,
		player.OverallGems,
		player.OverallEnemiesKilled,
		player.OverallDaggersHit,
		player.OverallDaggersFired,
		player.OverallAccuracy,
	)
	if err != nil {
		return err
	}
	return nil
}