package postgres

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB                     *sqlx.DB
	Games                  *GameModel
	States                 *StateModel
	Players                *PlayerModel
	ReplayPlayers          *ReplayPlayerModel
	SubmittedGames         *SubmittedGameModel
	MOTD                   *MOTDModel
	DiscordUsers           *DiscordUserModel
	Releases               *ReleaseModel
	Spawnsets              *SpawnsetModel
	News                   *NewsModel
	CollectorRuns          *CollectorRunModel
	CollectorPlayers       *CollectorPlayerModel
	CollectorHighScores    *CollectorHighScoreModel
	CollectorActivePlayers *CollectorActivePlayerModel
	CollectorNewPlayers    *CollectorNewPlayerModel
}

func NewPostgres(client *http.Client, db *sqlx.DB) *Postgres {
	return &Postgres{
		DB:                     db,
		Games:                  &GameModel{DB: db},
		States:                 &StateModel{DB: db},
		Players:                &PlayerModel{DB: db},
		ReplayPlayers:          &ReplayPlayerModel{DB: db},
		SubmittedGames:         &SubmittedGameModel{DB: db, Client: client},
		MOTD:                   &MOTDModel{DB: db},
		DiscordUsers:           &DiscordUserModel{DB: db},
		Releases:               &ReleaseModel{DB: db},
		Spawnsets:              &SpawnsetModel{DB: db},
		News:                   &NewsModel{DB: db},
		CollectorRuns:          &CollectorRunModel{DB: db},
		CollectorPlayers:       &CollectorPlayerModel{DB: db},
		CollectorHighScores:    &CollectorHighScoreModel{DB: db},
		CollectorActivePlayers: &CollectorActivePlayerModel{DB: db},
		CollectorNewPlayers:    &CollectorNewPlayerModel{DB: db},
	}
}
