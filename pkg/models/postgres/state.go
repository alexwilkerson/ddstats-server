package postgres

import (
	"database/sql"
)

// StateModel wraps database connection
type StateModel struct {
	DB *sql.DB
}
