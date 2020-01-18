package postgres

import (
	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type CollectorRunModel struct {
	DB *sqlx.DB
}

func (crm *CollectorRunModel) CreateNew(tx *sqlx.Tx) (int, error) {
	var id int
	stmt := `
		INSERT INTO collector_run DEFAULT VALUES returning id`
	err := tx.Get(&id, stmt)
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
		WHERE run_time != 0
		ORDER BY id DESC LIMIT 1`
	err := crm.DB.Get(&cr, stmt)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

func (crm *CollectorRunModel) Update(tx *sqlx.Tx, cr *models.CollectorRun) error {
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
			global_default_daggers=$24,
			global_bronze_daggers=$25,
			global_silver_daggers=$26,
			global_gold_daggers=$27,
			global_devil_daggers=$28,
			since_game_time=$29,
			since_deaths=$30,
			since_gems=$31,
			since_enemies_killed=$32,
			since_daggers_hit=$33,
			since_daggers_fired=$34,
			since_accuracy=$35,
			since_bronze_daggers=$36,
			since_silver_daggers=$37,
			since_gold_daggers=$38,
			since_devil_daggers=$39,
			fallen=$40,
			swarmed=$41,
			impaled=$42,
			gored=$43,
			infested=$44,
			opened=$45,
			purged=$46,
			desecrated=$47,
			sacrificed=$48,
			eviscerated=$49,
			annihilated=$50,
			intoxicated=$51,
			envenmonated=$52,
			incarnated=$53,
			discarnated=$54,
			barbed=$55
		WHERE id=$56`
	_, err := tx.Exec(stmt,
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
		cr.GlobalDefaultDaggers,
		cr.GlobalBronzeDaggers,
		cr.GlobalSilverDaggers,
		cr.GlobalGoldDaggers,
		cr.GlobalDevilDaggers,
		cr.SinceGameTime,
		cr.SinceDeaths,
		cr.SinceGems,
		cr.SinceEnemiesKilled,
		cr.SinceDaggersHit,
		cr.SinceDaggersFired,
		cr.SinceAccuracy,
		cr.SinceBronzeDaggers,
		cr.SinceSilverDaggers,
		cr.SinceGoldDaggers,
		cr.SinceDevilDaggers,
		cr.Fallen,
		cr.Swarmed,
		cr.Impaled,
		cr.Gored,
		cr.Infested,
		cr.Opened,
		cr.Purged,
		cr.Desecrated,
		cr.Sacrificed,
		cr.Eviscerated,
		cr.Annihilated,
		cr.Intoxicated,
		cr.Envenmonated,
		cr.Incarnated,
		cr.Discarnated,
		cr.Barbed,
		cr.ID,
	)
	return err
}

func (crm *CollectorRunModel) InsertNew(tx *sqlx.Tx) (int, error) {
	var id int
	stmt := `
		INSERT INTO collector_run
		DEFAULT VALUES
		RETURNING id`
	err := tx.Get(&id, stmt)
	if err != nil {
		return 0, err
	}
	return id, nil
}
