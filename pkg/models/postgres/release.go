package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
	"github.com/jmoiron/sqlx"
)

type ReleaseModel struct {
	DB *sqlx.DB
}

func (rm *ReleaseModel) Select(version string) (*models.Release, error) {
	var release models.Release
	stmt := `
		SELECT *
		FROM release
		WHERE version=$1`
	err := rm.DB.Get(&release, stmt, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
	}
	return &release, nil
}

func (rm *ReleaseModel) GetAll(pageSize, pageNum int) ([]*models.Release, error) {
	var releases []*models.Release
	stmt := fmt.Sprintf(`
		SELECT *
		FROM release
		ORDER BY time_stamp DESC LIMIT %d OFFSET %d`, pageSize, pageNum-1)
	err := rm.DB.Select(&releases, stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
	}
	return releases, nil
}

func (rm *ReleaseModel) GetTotalCount() (int, error) {
	var count int
	stmt := `
		SELECT COUNT(1)
		FROM release`
	err := rm.DB.Get(&count, stmt)
	if err != nil {
		return 0, err
	}
	return count, nil
}
