package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

// GameModel wraps database connection
type GameModel struct {
	DB *sqlx.DB
}

const (
	v3SurvivalHashA  = "5ff43e37d0f85e068caab5457305754e"
	v3SurvivalHashB  = "569fead87abf4d30fdee4231a6398051"
	defaultSpawnset  = "v3"
	pacifistSpawnset = "pacifist"
)

// GetTop retrieves a slice of the top games in the database with a given limit
func (g *GameModel) GetTop(limit int) ([]*models.GameWithName, error) {
	games := []*models.GameWithName{}

	stmt := fmt.Sprintf(`
		SELECT
			game.id,
			player_id,
			p1.player_name,
			granularity,
			round(game.game_time, 4) as game_time,
			death_type.name as death_type,
			game.gems,
			game.homing_daggers,
			game.daggers_fired,
			game.daggers_hit,
			round(divzero(game.daggers_hit, game.daggers_fired)*100, 2) as accuracy,
			game.enemies_alive,
			game.enemies_killed,
			time_stamp,
			replay_player_id,
			CASE WHEN replay_player_id=0 THEN '' WHEN p2.id IS NULL THEN 'unknown' ELSE p2.player_name END AS replay_player_name,
			CASE WHEN spawnset.survival_hash IS NULL THEN 'unknown' ELSE spawnset.spawnset_name END AS spawnset,
			version,
			level_two_time,
			level_three_time,
			level_four_time,
			homing_daggers_max_time,
			enemies_alive_max_time,
			homing_daggers_max,
			enemies_alive_max
		FROM game JOIN player p1 ON game.player_id=p1.id JOIN death_type ON game.death_type=death_type.id
			NATURAL LEFT JOIN spawnset
			LEFT JOIN replay_player p2 ON game.replay_player_id=p2.id
		WHERE replay_player_id=0 AND (spawnset.spawnset_name='%s')
		ORDER BY game_time DESC LIMIT %d`, defaultSpawnset, limit)
	err := g.DB.Select(&games, stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return games, nil
}

// GetRecent retrieves a slice of users using a specified page size and page num starting at 1
func (g *GameModel) GetRecent(playerID, pageSize, pageNum int, sortBy, sortDir string) ([]*models.GameWithName, string, error) {
	var where string
	if playerID != 0 {
		where = fmt.Sprintf("WHERE game.player_id=$1 AND game.replay_player_id=0")
	}

	if sortBy == "" {
		sortBy = "id"
		sortDir = "desc"
	}

	games := []*models.GameWithName{}

	stmt := fmt.Sprintf(`
		SELECT
			game.id,
			player_id,
			p1.player_name,
			granularity,
			round(game.game_time, 4) as game_time,
			death_type.name as death_type,
			game.gems,
			game.homing_daggers,
			game.daggers_fired,
			game.daggers_hit,
			round(divzero(game.daggers_hit, game.daggers_fired)*100, 2) as accuracy,
			game.enemies_alive,
			game.enemies_killed,
			time_stamp,
			replay_player_id,
			CASE WHEN replay_player_id=0 THEN '' WHEN p2.id IS NULL THEN 'unknown' ELSE p2.player_name END AS replay_player_name,
			CASE WHEN spawnset.survival_hash IS NULL THEN 'unknown' ELSE spawnset.spawnset_name END AS spawnset,
			version,
			level_two_time,
			level_three_time,
			level_four_time,
			homing_daggers_max_time,
			enemies_alive_max_time,
			homing_daggers_max,
			enemies_alive_max
		FROM game JOIN player p1 ON game.player_id=p1.id JOIN death_type ON game.death_type=death_type.id
			NATURAL LEFT JOIN spawnset
			LEFT JOIN replay_player p2 ON game.replay_player_id=p2.id %s
		ORDER BY %s %s LIMIT %d OFFSET %d`, where, sortBy, sortDir, pageSize, (pageNum-1)*pageSize)
	var err error
	if playerID != 0 {
		err = g.DB.Select(&games, stmt, playerID)
	} else {
		err = g.DB.Select(&games, stmt)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", models.ErrNoRecord
		}
		return nil, "", err
	}

	if playerID != 0 && len(games) > 0 {
		return games, games[0].PlayerName, nil
	}

	return games, "", nil
}

// GetLeaderboardPaginated is a function
func (g *GameModel) GetLeaderboardPaginated(spawnset string, pageSize, pageNum int) ([]*models.GameWithName, error) {
	games := []*models.GameWithName{}

	var where string
	var enemies string

	if spawnset == pacifistSpawnset {
		where = `
			WHERE spawnset_name='v3'
			AND (replay_player_id=0 OR replay_player_id=player_id)
			AND enemies_killed=0
			AND daggers_hit=0
			AND homing_daggers=0
			AND game_time < 300`
		enemies = "AND game.enemies_killed=0"
	} else {
		where = "WHERE spawnset_name=$1 AND (replay_player_id=0 OR replay_player_id=player_id)"
	}

	stmt := fmt.Sprintf(`
		WITH max_game AS (
			SELECT
				id,
				game.player_id,
				granularity,
				round(game_time, 4) AS game_time,
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
			FROM game INNER JOIN (
				SELECT DISTINCT ON (player_id) player_id, round(MAX(game_time), 4) AS max_game_time
				FROM game
				NATURAL LEFT JOIN spawnset
				%s
				GROUP BY player_id) gg ON game.player_id=gg.player_id AND round(game.game_time, 4)=gg.max_game_time %s),
		min_replay AS(
			SELECT player_id, MIN(replay_player_id) AS min_replay 
			FROM max_game
			group by player_id
		)

		SELECT ROW_NUMBER() OVER (ORDER BY ggg.game_time DESC) AS rank, ggg.* FROM (
			SELECT DISTINCT ON (player_id, game_time)
				max_game.id,
				p1.player_name,
				max_game.player_id,
				max_game.granularity,
				max_game.game_time,
				death_type.name AS death_type,
				max_game.gems,
				max_game.homing_daggers,
				max_game.daggers_fired,
				max_game.daggers_hit,
				round(divzero(max_game.daggers_hit, max_game.daggers_fired)*100, 2) as accuracy,
				max_game.enemies_alive,
				max_game.enemies_killed,
				max_game.replay_player_id,
				max_game.time_stamp,
				CASE WHEN spawnset.survival_hash IS NULL THEN 'unknown' ELSE spawnset.spawnset_name END AS spawnset,
				max_game.version,
				max_game.level_two_time,
				max_game.level_three_time,
				max_game.level_four_time,
				max_game.homing_daggers_max_time, 
				max_game.enemies_alive_max_time,
				max_game.homing_daggers_max,
				max_game.enemies_alive_max
			FROM min_replay JOIN max_game
			ON min_replay.min_replay = max_game.replay_player_id AND min_replay.player_id = max_game.player_id
			NATURAL LEFT JOIN spawnset
			JOIN player p1 ON max_game.player_id=p1.id JOIN death_type ON max_game.death_type=death_type.id
			ORDER BY game_time DESC LIMIT %d OFFSET %d
		) ggg`, where, enemies, pageSize, (pageNum-1)*pageSize)
	var err error
	if spawnset == pacifistSpawnset {
		err = g.DB.Select(&games, stmt)
	} else {
		err = g.DB.Select(&games, stmt, spawnset)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return games, nil
}

// GetLeaderboard is a function
func (g *GameModel) GetLeaderboard(spawnset string) ([]*models.GameWithName, error) {
	games := []*models.GameWithName{}

	var where string
	var enemies string

	if spawnset == pacifistSpawnset {
		where = `
			WHERE spawnset_name='v3'
			AND (replay_player_id=0 OR replay_player_id=player_id)
			AND enemies_killed=0
			AND daggers_hit=0
			AND homing_daggers=0
			AND game_time < 300`
		enemies = "AND game.enemies_killed=0"
	} else {
		where = "WHERE spawnset_name=$1 AND (replay_player_id=0 OR replay_player_id=player_id)"
	}

	stmt := fmt.Sprintf(`
		WITH max_game AS (
			SELECT
				id,
				game.player_id,
				granularity,
				round(game_time, 4) AS game_time,
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
			FROM game INNER JOIN (
				SELECT DISTINCT ON (player_id) player_id, round(MAX(game_time), 4) AS max_game_time
				FROM game
				NATURAL LEFT JOIN spawnset
				%s
				GROUP BY player_id) gg ON game.player_id=gg.player_id AND round(game.game_time, 4)=gg.max_game_time %s),
		min_replay AS(
			SELECT player_id, MIN(replay_player_id) AS min_replay 
			FROM max_game
			group by player_id
		)

		SELECT ROW_NUMBER() OVER (ORDER BY ggg.game_time DESC) AS rank, ggg.* FROM (
			SELECT DISTINCT ON (player_id, game_time)
				max_game.id,
				p1.player_name,
				max_game.player_id,
				max_game.granularity,
				max_game.game_time,
				death_type.name AS death_type,
				max_game.gems,
				max_game.homing_daggers,
				max_game.daggers_fired,
				max_game.daggers_hit,
				round(divzero(max_game.daggers_hit, max_game.daggers_fired)*100, 2) as accuracy,
				max_game.enemies_alive,
				max_game.enemies_killed,
				max_game.replay_player_id,
				max_game.time_stamp,
				CASE WHEN spawnset.survival_hash IS NULL THEN 'unknown' ELSE spawnset.spawnset_name END AS spawnset,
				max_game.version,
				max_game.level_two_time,
				max_game.level_three_time,
				max_game.level_four_time,
				max_game.homing_daggers_max_time, 
				max_game.enemies_alive_max_time,
				max_game.homing_daggers_max,
				max_game.enemies_alive_max
			FROM min_replay JOIN max_game
			ON min_replay.min_replay = max_game.replay_player_id AND min_replay.player_id = max_game.player_id
			NATURAL LEFT JOIN spawnset
			JOIN player p1 ON max_game.player_id=p1.id JOIN death_type ON max_game.death_type=death_type.id
			ORDER BY game_time
		) ggg`, where, enemies)
	var err error
	if spawnset == pacifistSpawnset {
		err = g.DB.Select(&games, stmt)
	} else {
		err = g.DB.Select(&games, stmt, spawnset)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return games, nil
}

// GetLeaderboardTotalCount returns the total number of games in the for leaderboards
func (g *GameModel) GetLeaderboardTotalCount(spawnset string) (int, error) {
	var err error
	var gameCount int

	var stmt string

	if spawnset == pacifistSpawnset {
		stmt = `
		SELECT COUNT(1) FROM (
			SELECT MAX(game_time) AS max_game_time
			FROM game
			NATURAL LEFT JOIN spawnset
			WHERE spawnset_name='v3'
				AND (replay_player_id=0 OR replay_player_id=player_id)
				AND enemies_killed=0
				AND daggers_hit=0
				AND homing_daggers=0
				AND game_time < 300
			GROUP BY player_id
		) g`
	} else {
		stmt = `
		SELECT COUNT(1) FROM (
			SELECT MAX(game_time) AS max_game_time
			FROM game
			NATURAL LEFT JOIN spawnset
			WHERE spawnset_name=$1 AND (replay_player_id=0 OR replay_player_id=player_id)
			GROUP BY player_id
		) g`
	}

	if spawnset == pacifistSpawnset {
		err = g.DB.QueryRow(stmt).Scan(&gameCount)
	} else {
		err = g.DB.QueryRow(stmt, spawnset).Scan(&gameCount)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}
	return gameCount, nil
}

// Get retrieves the entire game object
func (g *GameModel) Get(id int) (*models.GameWithName, error) {
	var game models.GameWithName
	stmt := `
		SELECT
			game.id,
			player_id,
			p1.player_name,
			granularity,
			round(game.game_time, 4) as game_time,
			death_type.name as death_type,
			game.gems,
			game.homing_daggers,
			game.daggers_fired,
			game.daggers_hit,
			round(divzero(game.daggers_hit, game.daggers_fired)*100, 2) as accuracy,
			game.enemies_alive,
			game.enemies_killed,
			time_stamp,
			replay_player_id,
			CASE WHEN replay_player_id=0 THEN '' WHEN p2.id IS NULL THEN 'unknown' ELSE p2.player_name END AS replay_player_name,
			CASE WHEN spawnset.survival_hash IS NULL THEN 'unknown' ELSE spawnset.spawnset_name END AS spawnset,
			version,
			level_two_time,
			level_three_time,
			level_four_time,
			homing_daggers_max_time,
			enemies_alive_max_time,
			homing_daggers_max,
			enemies_alive_max
		FROM game JOIN player p1 ON game.player_id=p1.id JOIN death_type ON game.death_type=death_type.id
			NATURAL LEFT JOIN spawnset
			LEFT JOIN replay_player p2 ON game.replay_player_id=p2.id
		WHERE game.id=$1`
	err := g.DB.Get(&game, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return &game, nil
}

// GetTotalCount returns the total number of games in the database
func (g *GameModel) GetTotalCount(playerID int) (int, error) {
	var err error
	var stmt string
	var gameCount int
	if playerID != 0 {
		stmt = `
			SELECT COUNT(1)
			FROM game
			WHERE player_id=$1 AND replay_player_id=0`
		err = g.DB.QueryRow(stmt, playerID).Scan(&gameCount)
	} else {
		stmt = `
			SELECT COUNT(1)
			FROM game`
		err = g.DB.QueryRow(stmt).Scan(&gameCount)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}
	return gameCount, nil
}
