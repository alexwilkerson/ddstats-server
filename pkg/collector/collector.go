package collector

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/models"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
)

const (
	maxLimit int = 1000
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
		c.Stop()
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
	runID, err := c.DB.CollectorRuns.CreateNew()
	if err != nil {
		c.errorLog.Printf("collector error: %v", err)
		return
	}
	run := models.CollectorRun{ID: runID}
	for ok := true; ok; ok = leaderboard.PlayerCount > 0 {
		select {
		case <-c.quit:
			return
		default:
			leaderboard, err = c.DDAPI.GetLeaderboard(maxLimit, offset)
			if err != nil {
				c.errorLog.Printf("collector error: %v", err)
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
					return
				default:
					previousPlayer, err := c.DB.CollectorPlayers.Select(int(player.PlayerID))
					if err != nil && !errors.Is(err, models.ErrNoRecord) {
						c.errorLog.Printf("collector error: %v", err)
						return
					}
					if errors.Is(err, models.ErrNoRecord) {
						c.calculateNewPlayer(player)
					} else {
						c.calculatePlayer(player, previousPlayer)
					}
					err = c.DB.CollectorPlayers.UpsertPlayer(player, run.ID)
					if err != nil {
						c.infoLog.Printf("%+v", player)
						c.errorLog.Printf("collector error: %v", err)
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
	err = c.DB.CollectorRuns.Update(&run)
	if err != nil {
		c.errorLog.Printf("collector error: %v", err)
	}
}

func (c *Collector) Stop() {
	close(c.quit)
	<-c.done
	c.infoLog.Printf("%d Players recorded to database...", c.totalPlayers)
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

func (c *Collector) calculatePlayer(fromDDAPI *ddapi.Player, fromDB *models.CollectorPlayer) {
	overallDeaths := int(fromDDAPI.OverallDeaths) - fromDB.OverallDeaths
	if overallDeaths < 1 {
		return
	}
	c.activePlayers++
	c.playerDeaths += overallDeaths
	gameTime := float64(fromDDAPI.GameTime) - fromDB.GameTime
	if gameTime > 0 {
		fmt.Println("old rank:", fromDB.GameTime, ";new rank:", fromDDAPI.GameTime)
		fmt.Println("gameTime", gameTime)
		c.playersWithNewScores++
		c.playerImprovementTime += gameTime
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

func (c *Collector) calculateNewPlayer(p *ddapi.Player) {
	overallDeaths := int(p.OverallDeaths)
	if overallDeaths < 1 {
		return
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
}
