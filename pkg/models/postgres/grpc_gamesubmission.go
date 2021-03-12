package postgres

import (
	"database/sql"
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

// CheckDuplicate takes a submitted game, checks if it's a replay...
// if it is, it will try to find a matching game recording to make sure
// the game hasn't already been previously recorded into the database.
// This is to make sure people who are watching replays of their own games
// won't record the same games over and over and over again.
func (gsm *GameSubmissionModel) CheckDuplicate(game *pb.SubmitGameRequest) (bool, int32, error) {
	var id int

	stmt := `
		SELECT id
		FROM game
		WHERE
			player_id=$1 AND
			game_time=$2 AND
			death_type=$3 AND
			gems=$4 AND
			homing_daggers=$5 AND
			daggers_fired=$6 AND
			daggers_hit=$7 AND
			enemies_alive=$8 AND
			enemies_killed=$9 AND
			homing_daggers_max=$10 AND
			enemies_alive_max=$11 AND
			(replay_player_id=0 OR replay_player_id=$12) AND
			total_gems=$13 AND
			gems_despawned=$14 AND
			gems_eaten=$15 AND
			daggers_eaten=$16
		ORDER BY id ASC
		LIMIT 1`
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
		game.HomingDaggersMax,
		game.EnemiesAliveMax,
		game.ReplayPlayerID,
		game.TotalGems,
		game.GemsDespawned,
		game.GemsEaten,
		game.DaggersEaten,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, 0, nil
		}
		return false, 0, err
	}
	return true, int32(id), nil
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
			is_replay,
			per_enemy_alive_count,
			per_enemy_kill_count)
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
		game.TimeLvl2,
		game.TimeLvl3,
		game.TimeLvl4,
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
