package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-api/pkg/models"
)

//GameModel wraps database connection
type GameModel struct {
	DB *sql.DB
}

//Get retrives the entire game obeject
func (g *GameModel) Get(id int) (*models.Game, error) {

	stmt := `SELECT * FROM game WHERE id=$1`
	row := g.DB.QueryRow(stmt, id)
	//This will hold the values of the retreived record
	gameModel := &models.Game{}
	err := row.Scan(&gameModel.ID,
		&gameModel.PlayerID,
		&gameModel.Granularity,
		&gameModel.GameTime,
		&gameModel.DeathType,
		&gameModel.Gems,
		&gameModel.HomingDaggers,
		&gameModel.DaggersFired,
		&gameModel.DaggersHit,
		&gameModel.EnemiesAlive,
		&gameModel.EnemiesKilled,
		&gameModel.TimeStamp,
		&gameModel.ReplayPlayerID,
		&gameModel.SurvivalHash,
		&gameModel.Version,
		&gameModel.LevelTwoTime,
		&gameModel.LevelThreeTime,
		&gameModel.LevelFourTime,
		&gameModel.HomingDaggersMaxTime,
		&gameModel.EnemiesAliveMaxTime,
		&gameModel.HomingDaggersMax,
		&gameModel.EnemiesAliveMax)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return gameModel, nil
}

//GetGems returns how many Gems in the game
func (g *GameModel) GetGems(id int) (int, error) {
	return 0, nil
}

//GetHomingDaggers returns how many homing daggers
func (g *GameModel) GetHomingDaggers(id int) (int, error) {
	return 0, nil
}

//GetAccuracy returns the game total accuracy
func (g *GameModel) GetAccuracy(id int) (int, error) {
	return 0, nil
}

//GetEnemiesAlive returns how many enemies are still alive
func (g *GameModel) GetEnemiesAlive(id int) (int, error) {
	return 0, nil
}

//GetEnemiesKilled return how many enemies had been killed
func (g *GameModel) GetEnemiesKilled(id int) (int, error) {
	return 0, nil
}
