package postgres

import "github.com/jmoiron/sqlx"

type CollectorRunModel struct {
	DB *sqlx.DB
}

func (crm *CollectorRunModel) SelectLastRunID() (int, error) {
	var id int
	stmt := `
		SELECT id
		FROM collector_run
		ORDER BY id DESC LIMIT 1`
	err := crm.DB.Get(&id, stmt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (crm *CollectorRunModel) InsertNew() (int, error) {
	var id int
	stmt := `
		INSERT INTO collector_run DEFAULT VALUES returning id`
	err := crm.DB.Get(&id, stmt)
	if err != nil {
		return 0, err
	}
	return id, nil
}
