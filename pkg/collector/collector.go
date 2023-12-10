package collector

import (
	"errors"
	"log"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
)

const (
	maxLimit int = 1000
)

const (
	BronzeDaggerThreshold float64 = 60
	SilverDaggerThreshold float64 = 120
	GoldDaggerThreshold   float64 = 250
	DevilDaggerThreshold  float64 = 500
)

const (
	Fallen       = "FALLEN"
	Swarmed      = "SWARMED"
	Impaled      = "IMPALED"
	Gored        = "GORED"
	Infested     = "INFESTED"
	Opened       = "OPENED"
	Purged       = "PURGED"
	Desecrated   = "DESECRATED"
	Sacrificed   = "SACRIFICED"
	Eviscerated  = "EVISCERATED"
	Annihilated  = "ANNIHILATED"
	Intoxicated  = "INTOXICATED"
	Envenmonated = "ENVENMONATED"
	Incarnated   = "INCARNATED"
	Discarnated  = "DISCARNATED"
	Barbed       = "BARBED"
)

type Collector struct {
	DDAPI                 *ddapi.API
	DB                    *postgres.Postgres
	infoLog               *log.Logger
	errorLog              *log.Logger
	totalPlayers          int
	activePlayers         int
	playersWithNewScores  int
	playersWithNewRanks   int
	playerImprovementTime float64
	playerRankImprovement int
	playerGameTime        float64
	playerDeaths          int
	playerGems            int
	playerEnemiesKilled   int
	playerDaggersHit      int
	playerDaggersFired    int
	quit                  chan struct{}
	done                  chan struct{}
}

func NewCollector(ddAPI *ddapi.API, db *postgres.Postgres, infoLog, errorLog *log.Logger) *Collector {
	return &Collector{
		DDAPI:    ddAPI,
		DB:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (c *Collector) Start() {
	defer func() {
		close(c.done)
	}()
	start := time.Now()
	var leaderboard *ddapi.Leaderboard
	var err error
	var offset int
	previousRun, err := c.DB.CollectorRuns.SelectLastRunID()
	if err != nil {
		c.errorLog.Printf("collector error: %v", err)
		return
	}
	tx, err := c.DB.DB.Beginx()
	if err != nil {
		c.errorLog.Printf("collector error: %v", err)
		return
	}
	runID, err := c.DB.CollectorRuns.CreateNew(tx)
	if err != nil {
		c.rollbackAndLogError(tx, err)
		return
	}
	run := models.CollectorRun{ID: runID}
	for ok := true; ok; ok = leaderboard.PlayerCount > 0 {
		select {
		case <-c.quit:
			err = tx.Rollback()
			if err != nil {
				c.errorLog.Printf("collector rollback error: %v", err)
			}
			return
		default:
			leaderboard, err = c.DDAPI.GetLeaderboard(maxLimit, offset)
			if err != nil {
				c.rollbackAndLogError(tx, err)
				return
			}
			// only run once
			if offset == 0 {
				initRun(&run, previousRun, leaderboard)
			}
			c.infoLog.Printf("collector offset: %d", offset)
			for _, player := range leaderboard.Players {
				select {
				case <-c.quit:
					c.infoLog.Println("Collector exiting prematurely. Rolling back database changes...")
					err = tx.Rollback()
					if err != nil {
						c.errorLog.Printf("collector rollback error: %v", err)
					}
					return
				default:
					err := c.DB.ReplayPlayers.Upsert(int(player.PlayerID), player.PlayerName)
					if err != nil {
						c.rollbackAndLogError(tx, err)
						return
					}
					switch player.DeathType {
					case Fallen:
						run.Fallen++
					case Swarmed:
						run.Swarmed++
					case Impaled:
						run.Impaled++
					case Gored:
						run.Gored++
					case Infested:
						run.Infested++
					case Opened:
						run.Opened++
					case Purged:
						run.Purged++
					case Desecrated:
						run.Desecrated++
					case Sacrificed:
						run.Sacrificed++
					case Eviscerated:
						run.Eviscerated++
					case Annihilated:
						run.Annihilated++
					case Intoxicated:
						run.Intoxicated++
					case Envenmonated:
						run.Envenmonated++
					case Incarnated:
						run.Incarnated++
					case Discarnated:
						run.Discarnated++
					case Barbed:
						run.Barbed++
					}
					previousPlayer, err := c.DB.CollectorPlayers.Select(int(player.PlayerID))
					if err != nil && !errors.Is(err, models.ErrNoRecord) {
						err = c.calculateNewPlayer(tx, &run, player)
						return
					}
					var activePlayer bool
					if errors.Is(err, models.ErrNoRecord) {
						err = c.calculateNewPlayer(tx, &run, player)
						activePlayer = true
						if err != nil {
							c.rollbackAndLogError(tx, err)
							return
						}
					} else {
						activePlayer, err = c.calculatePlayer(tx, &run, player, previousPlayer)
						if err != nil {
							c.rollbackAndLogError(tx, err)
							return
						}
					}
					if activePlayer {
						lastActive := time.Now()
						err = c.DB.CollectorPlayers.UpsertPlayer(tx, player, run.ID, &lastActive)
					} else {
						err = c.DB.CollectorPlayers.UpsertPlayer(tx, player, run.ID, nil)
					}
					if err != nil {
						c.rollbackAndLogError(tx, err)
						return
					}
					c.totalPlayers++
				}
			}
			offset += leaderboard.PlayerCount
		}
	}
	c.compileRunStats(&run, previousRun)
	run.RunTime = models.Duration(time.Since(start))
	err = c.DB.CollectorRuns.Update(tx, &run)
	if err != nil {
		c.rollbackAndLogError(tx, err)
		return
	}
	err = tx.Commit()
	if err != nil {
		c.errorLog.Printf("collector commit error: %v", err)
	}
	c.infoLog.Printf("%d Players recorded to database...", c.totalPlayers)
}

func (c *Collector) Stop() {
	close(c.quit)
	<-c.done
}

func (c *Collector) Done() chan struct{} {
	return c.done
}

func (c *Collector) rollbackAndLogError(tx *sqlx.Tx, err error) {
	c.errorLog.Printf("collector error: %v", err)
	err = tx.Rollback()
	if err != nil {
		c.errorLog.Printf("collector rollback error: %v", err)
	}
}

func initRun(run *models.CollectorRun, previousRun *models.CollectorRun, leaderboard *ddapi.Leaderboard) {
	run.GlobalPlayers = int(leaderboard.GlobalPlayerCount)
	run.GlobalGameTime = float64(leaderboard.GlobalGameTime)
	run.GlobalDeaths = int(leaderboard.GlobalDeaths)
	run.GlobalGems = int(leaderboard.GlobalGems)
	run.GlobalEnemiesKilled = int(leaderboard.GlobalEnemiesKilled)
	run.GlobalDaggersHit = int(leaderboard.GlobalDaggersHit)
	run.GlobalDaggersFired = int(leaderboard.GlobalDaggersFired)
	run.GlobalAccuracy = float64(leaderboard.GlobalAccuracy)
	run.SinceGameTime = run.GlobalGameTime - previousRun.GlobalGameTime
	run.SinceDeaths = run.GlobalDeaths - previousRun.GlobalDeaths
	run.SinceGems = run.GlobalGems - previousRun.GlobalGems
	run.SinceEnemiesKilled = run.GlobalEnemiesKilled - previousRun.GlobalEnemiesKilled
	run.SinceDaggersHit = run.GlobalDaggersHit - previousRun.GlobalDaggersHit
	run.SinceDaggersFired = run.GlobalDaggersFired - previousRun.GlobalDaggersFired
	if run.SinceDaggersFired != 0 {
		run.SinceAccuracy = float64(run.SinceDaggersHit) / float64(run.SinceDaggersFired) * 100
	}
}

func (c *Collector) compileRunStats(run *models.CollectorRun, previousRun *models.CollectorRun) {
	run.NewPlayers = run.GlobalPlayers - previousRun.GlobalPlayers
	run.ActivePlayers = c.activePlayers
	run.InactivePlayers = run.GlobalPlayers - run.ActivePlayers
	run.PlayersWithNewScores = c.playersWithNewScores
	run.PlayersWithNewRanks = c.playersWithNewRanks
	if c.playersWithNewScores != 0 {
		run.AverageImprovementTime = c.playerImprovementTime / float64(c.playersWithNewScores)
	}
	if c.playersWithNewRanks != 0 {
		run.AverageRankImprovement = float64(c.playerRankImprovement) / float64(c.playersWithNewRanks)
	}
	if c.playerDeaths != 0 {
		run.AverageGameTimePerActivePlayer = c.playerGameTime / float64(c.playerDeaths)
	}
	activePlayers := float64(c.activePlayers)
	if activePlayers != 0 {
		run.AverageDeathsPerActivePlayer = float64(c.playerDeaths) / activePlayers
		run.AverageGemsPerActivePlayer = float64(c.playerGems) / activePlayers
		run.AverageEnemiesKilledPerActivePlayer = float64(c.playerEnemiesKilled) / activePlayers
		run.AverageDaggersHitPerActivePlayer = float64(c.playerDaggersHit) / activePlayers
		run.AverageDaggersFiredPerActivePlayer = float64(c.playerDaggersFired) / activePlayers
		if run.AverageDaggersFiredPerActivePlayer != 0 {
			run.AverageAccuracyPerActivePlayer = run.AverageDaggersHitPerActivePlayer / run.AverageDaggersFiredPerActivePlayer * 100
		}
	}
}

func (c *Collector) calculatePlayer(tx *sqlx.Tx, run *models.CollectorRun, fromDDAPI *ddapi.Player, fromDB *models.CollectorPlayer) (bool, error) {
	newDagger := calculateDaggers(run, fromDDAPI, fromDB)
	sinceDeaths := int(fromDDAPI.OverallDeaths) - fromDB.OverallDeaths
	sinceGameTime := float64(fromDDAPI.OverallGameTime) - fromDB.OverallGameTime
	if sinceDeaths < 1 || sinceGameTime == 0 {
		return false, nil
	}
	rankImprovement := fromDB.Rank - int(fromDDAPI.Rank)
	if rankImprovement > 0 {
		c.playersWithNewRanks++
		c.playerRankImprovement += rankImprovement
	} else {
		rankImprovement = 0
	}
	c.activePlayers++
	c.playerDeaths += sinceDeaths
	gameTimeImprovement := float64(fromDDAPI.GameTime) - fromDB.GameTime
	if gameTimeImprovement > 0 {
		c.playersWithNewScores++
		c.playerImprovementTime += gameTimeImprovement
		if newDagger {
			c.DB.CollectorHighScores.Insert(tx, run.ID, int(fromDDAPI.PlayerID), float64(fromDDAPI.GameTime))
		}
	}
	err := c.DB.CollectorActivePlayers.Insert(tx, run.ID, int(fromDDAPI.PlayerID), int(fromDDAPI.Rank), rankImprovement, float64(fromDDAPI.GameTime), gameTimeImprovement, sinceGameTime, sinceDeaths)
	if err != nil {
		return false, err
	}
	c.playerGameTime += sinceGameTime
	c.playerDeaths += int(fromDDAPI.OverallDeaths) - fromDB.OverallDeaths
	c.playerGems += int(fromDDAPI.OverallGems) - fromDB.OverallGems
	c.playerEnemiesKilled += int(fromDDAPI.OverallEnemiesKilled) - fromDB.OverallEnemiesKilled
	c.playerDaggersHit += int(fromDDAPI.OverallDaggersHit) - fromDB.OverallDaggersHit
	c.playerDaggersFired += int(fromDDAPI.OverallDaggersFired) - fromDB.OverallDaggersFired
	return true, nil
}

func (c *Collector) calculateNewPlayer(tx *sqlx.Tx, run *models.CollectorRun, p *ddapi.Player) error {
	calculateDaggers(run, p, nil)
	err := c.DB.CollectorPlayers.NewPlayer(tx, int(p.PlayerID))
	if err != nil {
		return err
	}
	err = c.DB.CollectorNewPlayers.Insert(tx, run.ID, int(p.PlayerID), int(p.Rank), float64(p.GameTime))
	if err != nil {
		return err
	}
	overallDeaths := int(p.OverallDeaths)
	if overallDeaths < 1 {
		return nil
	}
	err = c.DB.CollectorActivePlayers.Insert(tx, run.ID, int(p.PlayerID), int(p.Rank), 0, float64(p.GameTime), 0, float64(p.OverallGameTime), overallDeaths)
	if err != nil {
		return err
	}
	if p.GameTime >= BronzeDaggerThreshold {
		err = c.DB.CollectorHighScores.Insert(tx, run.ID, int(p.PlayerID), float64(p.GameTime))
		if err != nil {
			return err
		}
	}
	c.playerDeaths += overallDeaths
	c.activePlayers++
	gameTime := float64(p.GameTime)
	if gameTime > 0 {
		c.playersWithNewScores++
		c.playerImprovementTime += gameTime
	}
	c.playerGameTime += float64(p.OverallGameTime)
	c.playerGems += int(p.OverallGems)
	c.playerEnemiesKilled += int(p.OverallEnemiesKilled)
	c.playerDaggersHit += int(p.OverallDaggersHit)
	c.playerDaggersFired += int(p.OverallDaggersFired)
	return nil
}

func calculateDaggers(run *models.CollectorRun, fromDDAPI *ddapi.Player, fromDB *models.CollectorPlayer) bool {
	gameTimeFromDDAPI := float64(fromDDAPI.GameTime)
	switch {
	case gameTimeFromDDAPI >= DevilDaggerThreshold:
		run.GlobalDevilDaggers++
	case gameTimeFromDDAPI >= GoldDaggerThreshold:
		run.GlobalGoldDaggers++
	case gameTimeFromDDAPI >= SilverDaggerThreshold:
		run.GlobalSilverDaggers++
	case gameTimeFromDDAPI >= BronzeDaggerThreshold:
		run.GlobalBronzeDaggers++
	default:
		run.GlobalDefaultDaggers++
	}
	var newDagger bool
	if fromDB != nil {
		switch {
		case fromDB.GameTime < DevilDaggerThreshold && gameTimeFromDDAPI >= DevilDaggerThreshold:
			run.SinceDevilDaggers++
			newDagger = true
		case fromDB.GameTime < GoldDaggerThreshold && gameTimeFromDDAPI >= GoldDaggerThreshold:
			run.SinceGoldDaggers++
			newDagger = true
		case fromDB.GameTime < SilverDaggerThreshold && gameTimeFromDDAPI >= SilverDaggerThreshold:
			run.SinceSilverDaggers++
			newDagger = true
		case fromDB.GameTime < BronzeDaggerThreshold && gameTimeFromDDAPI >= BronzeDaggerThreshold:
			run.SinceBronzeDaggers++
			newDagger = true
		}
	} else { // if it's a new player, fromDB will be nil
		switch {
		case gameTimeFromDDAPI >= DevilDaggerThreshold:
			run.SinceDevilDaggers++
			newDagger = true
		case gameTimeFromDDAPI >= GoldDaggerThreshold:
			run.SinceGoldDaggers++
			newDagger = true
		case gameTimeFromDDAPI >= SilverDaggerThreshold:
			run.SinceSilverDaggers++
			newDagger = true
		case gameTimeFromDDAPI >= BronzeDaggerThreshold:
			run.SinceBronzeDaggers++
			newDagger = true
		}
	}
	return newDagger
}
