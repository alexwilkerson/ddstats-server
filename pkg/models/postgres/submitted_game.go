package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

type SubmittedGameModel struct {
	DB *sql.DB
}

func (sg *SubmittedGameModel) Insert(game *models.SubmittedGame) error {
	stmt := `INSERT INTO game(
				player_id,
				granularity,
				game_time,
				death_type,
				gems,
				homing_daggers,
				daggers_fired,
				daggers_hit,
				enemies_alive,
				enemies_killed,
				time_stamp,
				replay_player_id,
				survival_hash,
				version,
				level_two_time,
				level_three_time,
				level_four_time,
				homing_daggers_max_time,
				enemies_alive_max_time,
				homing_daggers_max,
				enemies_alive_max
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP,
				$11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
			RETURNING id`

	var gameID int
	err := sg.DB.QueryRow(stmt,
		game.PlayerID,
		game.Granularity,
		roundToNearest(game.GameTime, 4),
		game.DeathType,
		game.Gems,
		game.HomingDaggers,
		game.DaggersFired,
		game.DaggersHit,
		game.EnemiesAlive,
		game.EnemiesKilled,
		game.ReplayPlayerID,
		game.SurvivalHash,
		game.Version,
		game.Level2Time,
		game.Level3Time,
		game.Level4Time,
		game.HomingMaxTime,
		game.EnemiesAliveMaxTime,
		game.HomingMax,
		game.EnemiesAliveMax,
	).Scan(&gameID)
	if err != nil {
		return err
	}

	// Verify that all slices are of the same length
	if (len(game.GemsSlice)+
		len(game.HomingDaggersSlice)+
		len(game.DaggersHitSlice)+
		len(game.DaggersFiredSlice)+
		len(game.EnemiesAliveSlice)+
		len(game.EnemiesKilledSlice))/6 != len(game.GameTimeSlice) {
		return errors.New("invalid data")
	}

	states := StateModel{DB: sg.DB}

	for i := 0; i < len(game.GameTimeSlice); i++ {
		var accuracy float64
		if game.DaggersFired > 0 {
			accuracy = roundToNearest(float64(game.DaggersHit)/float64(game.DaggersFired)*100, 2)
		}
		state := models.State{
			GameID:        gameID,
			GameTime:      roundToNearest(game.GameTimeSlice[i], 4),
			Gems:          game.GemsSlice[i],
			HomingDaggers: game.HomingDaggersSlice[i],
			DaggersHit:    game.DaggersHitSlice[i],
			DaggersFired:  game.DaggersFiredSlice[i],
			Accuracy:      accuracy,
			EnemiesAlive:  game.EnemiesAliveSlice[i],
			EnemiesKilled: game.EnemiesKilledSlice[i],
		}
		err = states.Insert(&state)
		if err != nil {
			return err
		}
	}
	return nil
}
