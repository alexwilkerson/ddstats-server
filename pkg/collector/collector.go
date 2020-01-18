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
		c.errorLog.Printf("collector error: %v", err)
		err = tx.Rollback()
		if err != nil {
			c.errorLog.Printf("collector rollback error: %v", err)
		}
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
				c.errorLog.Printf("collector error: %v", err)
				err = tx.Rollback()
				if err != nil {
					c.errorLog.Printf("collector rollback error: %v", err)
				}
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
						c.errorLog.Printf("collector error: %v", err)
						err = tx.Rollback()
						if err != nil {
							c.errorLog.Printf("collector rollback error: %v", err)
						}
						return
					}
					if errors.Is(err, models.ErrNoRecord) {
						err = c.calculateNewPlayer(tx, runID, player)
						if err != nil {
							c.errorLog.Printf("collector error: %v", err)
							err = tx.Rollback()
							if err != nil {
								c.errorLog.Printf("collector rollback error: %v", err)
							}
							return
						}
					} else {
						c.calculatePlayer(tx, runID, player, previousPlayer)
					}
					err = c.DB.CollectorPlayers.UpsertPlayer(tx, player, run.ID)
					if err != nil {
						c.errorLog.Printf("collector error: %v", err)
						err = tx.Rollback()
						if err != nil {
							c.errorLog.Printf("collector rollback error: %v", err)
						}
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
		err = tx.Rollback()
		c.errorLog.Printf("collector error: %v", err)
		if err != nil {
			c.errorLog.Printf("collector rollback error: %v", err)
		}
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
	run.AverageGameTimePerActivePlayer = c.playerGameTime / float64(c.playerDeaths)
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

func (c *Collector) calculatePlayer(tx *sqlx.Tx, runID int, fromDDAPI *ddapi.Player, fromDB *models.CollectorPlayer) {
	overallDeaths := int(fromDDAPI.OverallDeaths) - fromDB.OverallDeaths
	if overallDeaths < 1 {
		return
	}
	c.activePlayers++
	c.playerDeaths += overallDeaths
	gameTime := float64(fromDDAPI.GameTime) - fromDB.GameTime
	if gameTime > 0 {
		c.playersWithNewScores++
		c.playerImprovementTime += gameTime
		if fromDB.GameTime < DevilDaggerThreshold && float64(fromDDAPI.GameTime) >= DevilDaggerThreshold ||
			fromDB.GameTime < GoldDaggerThreshold && float64(fromDDAPI.GameTime) >= GoldDaggerThreshold ||
			fromDB.GameTime < SilverDaggerThreshold && float64(fromDDAPI.GameTime) >= SilverDaggerThreshold ||
			fromDB.GameTime < BronzeDaggerThreshold && float64(fromDDAPI.GameTime) >= BronzeDaggerThreshold {
			c.DB.CollectorHighScores.Insert(tx, runID, int(fromDDAPI.PlayerID), float64(fromDDAPI.GameTime))
		}
	}
	rank := int(fromDDAPI.Rank) - fromDB.Rank
	if rank > 0 {
		c.playersWithNewRanks++
		c.playerRankImprovement += rank
	}
	c.playerGameTime += float64(fromDDAPI.OverallGameTime) - fromDB.OverallGameTime
	c.playerDeaths += int(fromDDAPI.OverallDeaths) - fromDB.OverallDeaths
	c.playerGems += int(fromDDAPI.OverallGems) - fromDB.OverallGems
	c.playerEnemiesKilled += int(fromDDAPI.OverallEnemiesKilled) - fromDB.OverallEnemiesKilled
	c.playerDaggersHit += int(fromDDAPI.OverallDaggersHit) - fromDB.OverallDaggersHit
	c.playerDaggersFired += int(fromDDAPI.OverallDaggersFired) - fromDB.OverallDaggersFired
}

func (c *Collector) calculateNewPlayer(tx *sqlx.Tx, runID int, p *ddapi.Player) error {
	overallDeaths := int(p.OverallDeaths)
	if overallDeaths < 1 {
		return nil
	}
	if p.GameTime >= BronzeDaggerThreshold {
		err := c.DB.CollectorPlayers.NewPlayer(tx, int(p.PlayerID))
		if err != nil {
			return err
		}
		c.DB.CollectorHighScores.Insert(tx, runID, int(p.PlayerID), float64(p.GameTime))
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
