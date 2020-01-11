package postgres

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	Games          *GameModel
	States         *StateModel
	Players        *PlayerModel
	SubmittedGames *SubmittedGameModel
	MOTD           *MOTDModel
}

func NewPostgres(client *http.Client, db *sqlx.DB) *Postgres {
	return &Postgres{
		Games:          &GameModel{DB: db},
		States:         &StateModel{DB: db},
		Players:        &PlayerModel{DB: db},
		SubmittedGames: &SubmittedGameModel{DB: db, Client: client},
		MOTD:           &MOTDModel{DB: db},
	}
}
