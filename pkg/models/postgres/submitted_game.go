package postgres

import (
	"errors"
	"math"
	"net/http"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
	"github.com/alexwilkerson/ddstats-api/pkg/models"
	"github.com/jmoiron/sqlx"
)

// SubmittedGameModel wraps the database connection and http client
type SubmittedGameModel struct {
	DB     *sqlx.DB
	Client *http.Client
}

// Insert takes a submitted game and inserts the data into the game table,
// then iterates over all of the states and inserts each state into the state table
func (sg *SubmittedGameModel) Insert(game *models.SubmittedGame) (int, error) {
	// fixes possible older versions of client submitting
	if game.SurvivalHash == "" {
		game.SurvivalHash = "5ff43e37d0f85e068caab5457305754e"
	}

	stmt := `
		INSERT INTO game(
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
			enemies_alive_max)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id`
	var gameID int
	err := sg.DB.QueryRow(stmt,
		game.PlayerID,
		game.Granularity,
		game.GameTime,
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
		game.LevelTwoTime,
		game.LevelThreeTime,
		game.LevelFourTime,
		game.HomingMaxTime,
		game.EnemiesAliveMaxTime,
		game.HomingMax,
		game.EnemiesAliveMax,
	).Scan(&gameID)
	if err != nil {
		return 0, err
	}

	// Verify that all slices are of the same length
	if (len(game.GemsSlice)+
		len(game.HomingDaggersSlice)+
		len(game.DaggersHitSlice)+
		len(game.DaggersFiredSlice)+
		len(game.EnemiesAliveSlice)+
		len(game.EnemiesKilledSlice))/6 != len(game.GameTimeSlice) {
		return 0, errors.New("invalid data")
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
			return 0, err
		}
	}

	players := PlayerModel{sg.DB}
	_, err = players.Get(game.PlayerID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			ddAPI := ddapi.API{Client: sg.Client}
			player, err := ddAPI.UserByID(game.PlayerID)
			if err != nil {
				return 0, err
			}
			err = players.UpsertDDPlayer(player)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return gameID, nil
}

func roundToNearest(f float64, numberOfDecimalPlaces int) float64 {
	multiplier := math.Pow10(numberOfDecimalPlaces)
	return math.Round(f*multiplier) / multiplier
}
