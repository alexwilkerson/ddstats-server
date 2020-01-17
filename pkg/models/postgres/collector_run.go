package postgres

import (
	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type CollectorRunModel struct {
	DB *sqlx.DB
}

func (crm *CollectorRunModel) CreateNew() (int, error) {
	var id int
	stmt := `
		INSERT INTO collector_run DEFAULT VALUES returning id`
	err := crm.DB.Get(&id, stmt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (crm *CollectorRunModel) SelectLastRunID() (*models.CollectorRun, error) {
	var cr models.CollectorRun
	stmt := `
		SELECT *
		FROM collector_run
		ORDER BY id DESC LIMIT 1`
	err := crm.DB.Get(&cr, stmt)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

func (crm *CollectorRunModel) Update(cr *models.CollectorRun) error {
	stmt := `
		UPDATE collector_run
		SET
			run_time=$1,
			global_players=$2,
			new_players=$3,
			active_players=$4,
			inactive_players=$5,
			players_with_new_scores=$6,
			players_with_new_ranks=$7,
			average_improvement_time=$8,
			average_rank_improvement=$9,
			average_game_time_per_active_player=$10,
			average_deaths_per_active_player=$11,
			average_gems_per_active_player=$12,
			average_enemies_killed_per_active_player=$13,
			average_daggers_hit_per_active_player=$14,
			average_daggers_fired_per_active_player=$15,
			average_accuracy_per_active_player=$16,
			global_game_time=$17,
			global_deaths=$18,
			global_gems=$19,
			global_enemies_killed=$20,
			global_daggers_hit=$21,
			global_daggers_fired=$22,
			global_accuracy=$23,
			since_game_time=$24,
			since_deaths=$25,
			since_gems=$26,
			since_enemies_killed=$27,
			since_daggers_hit=$28,
			since_daggers_fired=$29,
			since_accuracy=$30
		WHERE id=$31`
	_, err := crm.DB.Exec(stmt,
		cr.RunTime,
		cr.GlobalPlayers,
		cr.NewPlayers,
		cr.ActivePlayers,
		cr.InactivePlayers,
		cr.PlayersWithNewScores,
		cr.PlayersWithNewRanks,
		cr.AverageImprovementTime,
		cr.AverageRankImprovement,
		cr.AverageGameTimePerActivePlayer,
		cr.AverageDeathsPerActivePlayer,
		cr.AverageGemsPerActivePlayer,
		cr.AverageEnemiesKilledPerActivePlayer,
		cr.AverageDaggersHitPerActivePlayer,
		cr.AverageDaggersFiredPerActivePlayer,
		cr.AverageAccuracyPerActivePlayer,
		cr.GlobalGameTime,
		cr.GlobalDeaths,
		cr.GlobalGems,
		cr.GlobalEnemiesKilled,
		cr.GlobalDaggersHit,
		cr.GlobalDaggersFired,
		cr.GlobalAccuracy,
		cr.SinceGameTime,
		cr.SinceDeaths,
		cr.SinceGems,
		cr.SinceEnemiesKilled,
		cr.SinceDaggersHit,
		cr.SinceDaggersFired,
		cr.SinceAccuracy,
		cr.ID,
	)
	return err
}

func (crm *CollectorRunModel) InsertNew() (int, error) {
	var id int
	stmt := `
		INSERT INTO collector_run
		DEFAULT VALUES
		RETURNING id`
	err := crm.DB.Get(&id, stmt)
	if err != nil {
		return 0, err
	}
	return id, nil
}
