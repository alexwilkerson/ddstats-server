package collector

import (
	"log"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
)

const (
	maxLimit int = 1000
)

type Collector struct {
	DDAPI       *ddapi.API
	DB          *postgres.Postgres
	infoLog     *log.Logger
	errorLog    *log.Logger
	playerCount int
	quit        chan struct{}
	done        chan struct{}
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
	var leaderboard *ddapi.Leaderboard
	var err error
	var offset int
	runID, err := c.DB.CollectorRuns.InsertNew()
	if err != nil {
		c.errorLog.Printf("collector error: %v", err)
		return
	}
	for ok := true; ok; ok = leaderboard.PlayerCount > 0 {
		select {
		case <-c.quit:
			close(c.done)
			return
		default:
			leaderboard, err = c.DDAPI.GetLeaderboard(maxLimit, offset)
			if err != nil {
				c.errorLog.Printf("collector error: %v", err)
				close(c.done)
				return
			}
			c.infoLog.Printf("collector offset: %d", offset)
			for _, player := range leaderboard.Players {
				select {
				case <-c.quit:
					close(c.done)
					return
				default:
					err = c.DB.CollectorPlayers.UpsertPlayer(player, runID)
					if err != nil {
						c.infoLog.Printf("%+v", player)
						c.errorLog.Printf("collector error: %v", err)
						close(c.done)
						return
					}
					c.playerCount++
				}
			}
			offset += leaderboard.PlayerCount
		}
	}
	close(c.done)
	c.Stop()
}

func (c *Collector) Stop() {
	close(c.quit)
	<-c.done
	c.infoLog.Printf("%d Players recorded to database...", c.playerCount)
}

func (c *Collector) Done() chan struct{} {
	return c.done
}
