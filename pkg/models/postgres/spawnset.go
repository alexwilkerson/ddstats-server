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
