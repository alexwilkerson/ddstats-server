package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
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

	var releaseNotes []models.ReleaseNote
	stmt = `
		SELECT *
		FROM release_note
		WHERE release_version=$1
		ORDER BY id ASC`
	err = rm.DB.Select(&releaseNotes, stmt, version)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	release.Notes = make([]string, 0, len(releaseNotes))
	for _, releaseNote := range releaseNotes {
		release.Notes = append(release.Notes, releaseNote.Note)
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

	for _, release := range releases {
		var releaseNotes []models.ReleaseNote
		stmt = `
			SELECT *
			FROM release_note
			WHERE release_version=$1
			ORDER BY id ASC`
		err = rm.DB.Select(&releaseNotes, stmt, release.Version)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		release.Notes = make([]string, 0, len(releaseNotes))
		for _, releaseNote := range releaseNotes {
			release.Notes = append(release.Notes, releaseNote.Note)
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
