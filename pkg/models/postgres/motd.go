package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
	"github.com/jmoiron/sqlx"
)

type MOTDModel struct {
	DB *sqlx.DB
}

func (m *MOTDModel) Get() (*models.MOTD, error) {
	var motd models.MOTD
	stmt := `
		SELECT *
		FROM message_of_the_day
		ORDER BY id DESC
		LIMIT 1`
	err := m.DB.Get(&motd, stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return &motd, nil
}
