package postgres

import (
	"errors"
	"net/http"

	pb "github.com/alexwilkerson/ddstats-server/gamesubmission"
	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"
	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type GameSubmissionModel struct {
	DB     *sqlx.DB
	Client *http.Client
}

func (gsm *GameSubmissionModel) Insert(game *pb.SubmitGameRequest) (int32, error) {
	stmt := `
		INSERT INTO game(
			player_id,
			game_time,
			death_type,
			gems,
			homing_daggers,
			daggers_fired,
			daggers_hit,
			enemies_alive,
			enemies_killed,
			replay_player_id,
			survival_hash,
			version,
			level_two_time,
			level_three_time,
			level_four_time,
			levi_down_time,
			orb_down_time,
			homing_daggers_max_time,
			enemies_alive_max_time,
			homing_daggers_max,
			enemies_alive_max,
			total_gems,
			level_gems,
			gems_despawned,
			gems_eaten,
			daggers_eaten,
			is_replay)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)
		RETURNING id`
	var gameID int
	err := gsm.DB.QueryRow(stmt,
		game.PlayerID,
		game.Time,
		game.DeathType,
		game.GemsCollected,
		game.HomingDaggers,
		game.DaggersFired,
		game.DaggersHit,
		game.EnemiesAlive,
		game.Kills,
		game.ReplayPlayerID,
		game.LevelHashMD5,
		game.Version,
		game.TimeLvll2,
		game.TimeLvll3,
		game.TimeLvll4,
		game.TimeLeviDown,
		game.TimeOrbDown,
		game.HomingDaggersMaxTime,
		game.EnemiesAliveMaxTime,
		game.HomingDaggersMax,
		game.EnemiesAliveMax,
		game.TotalGems,
		game.LevelGems,
		game.GemsDespawned,
		game.GemsEaten,
		game.DaggersEaten,
		game.IsReplay,
		pq.Array(game.PerEnemyAliveCount),
		pq.Array(game.PerEnemyKillcount),
	).Scan(&gameID)
	if err != nil {
		return 0, err
	}

	states := StateModel{DB: gsm.DB}
	for i, state := range game.Stats {
		var accuracy float64
		if state.DaggersFired > 0 {
			accuracy = roundToNearest(float64(state.DaggersHit)/float64(state.DaggersFired)*100, 2)
		}
		gameTime := float32(i)
		if i == len(game.Stats)-1 {
			gameTime = game.Time
		}
		state := models.State{
			GameID:             gameID,
			GameTime:           float64(gameTime),
			Gems:               int(state.GemsCollected),
			HomingDaggers:      int(state.HomingDaggers),
			DaggersHit:         int(state.DaggersHit),
			DaggersFired:       int(state.DaggersFired),
			Accuracy:           accuracy,
			EnemiesAlive:       int(state.EnemiesAlive),
			EnemiesKilled:      int(state.Kills),
			TotalGems:          state.TotalGems,
			LevelGems:          state.LevelGems,
			GemsDespawned:      state.GemsDespawned,
			GemsEaten:          state.GemsEaten,
			DaggersEaten:       state.DaggersEaten,
			PerEnemyAliveCount: state.PerEnemyAliveCount,
			PerEnemyKillCount:  state.PerEnemyKillCount,
		}
		err = states.InsertGRPC(&state)
		if err != nil {
			return 0, err
		}
	}

	players := PlayerModel{gsm.DB}
	_, err = players.Get(int(game.PlayerID))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			ddAPI := ddapi.API{Client: gsm.Client}
			player, err := ddAPI.UserByID(int(game.PlayerID))
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

	return int32(gameID), nil
}
