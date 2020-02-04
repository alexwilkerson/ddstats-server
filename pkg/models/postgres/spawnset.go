package postgres

import (
	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type SpawnsetModel struct {
	DB *sqlx.DB
}

func (sm *SpawnsetModel) SelectSpawnsetNames() ([]string, error) {
	var spawnsets []*models.Spawnset
	stmt := `
		SELECT DISTINCT spawnset_name
		FROM spawnset`
	err := sm.DB.Select(&spawnsets, stmt)
	if err != nil {
		return nil, err
	}
	spawnset_names := make([]string, 0, len(spawnsets))
	for _, spawnset := range spawnsets {
		spawnset_names = append(spawnset_names, spawnset.SpawnsetName)
	}
	return spawnset_names, nil
}

func (sm *SpawnsetModel) Select(name string) (*models.Spawnset, error) {
	var spawnset models.Spawnset
	stmt := `
		SELECT *
		FROM spawnset
		WHERE spawnset_name=$1 LIMIT 1`
	err := sm.DB.Get(&spawnset, stmt, name)
	if err != nil {
		return nil, err
	}
	return &spawnset, nil
}
