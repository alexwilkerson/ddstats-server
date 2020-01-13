package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
	"github.com/jmoiron/sqlx"
)

type NewsModel struct {
	DB *sqlx.DB
}

func (nm *NewsModel) Select(id int) (*models.News, error) {
	var news models.News
	stmt := `
		SELECT *
		FROM news
		WHERE id=$1`
	err := nm.DB.Get(&news, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
	}
	return &news, nil
}

func (nm *NewsModel) GetAll(pageSize, pageNum int) ([]*models.News, error) {
	var news []*models.News
	stmt := fmt.Sprintf(`
		SELECT *
		FROM news
		ORDER BY time_stamp DESC LIMIT %d OFFSET %d`, pageSize, pageNum-1)
	err := nm.DB.Select(&news, stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
	}
	return news, nil
}

func (nm *NewsModel) GetTotalCount() (int, error) {
	var count int
	stmt := `
		SELECT COUNT(1)
		FROM news`
	err := nm.DB.Get(&count, stmt)
	if err != nil {
		return 0, err
	}
	return count, nil
}
