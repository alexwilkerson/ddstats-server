package discord

import (
	"log"
	"strings"
	"sync"

	"github.com/alexwilkerson/ddstats-api/pkg/models/postgres"

	"github.com/alexwilkerson/ddstats-api/pkg/websocket"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"

	"github.com/bwmarrin/discordgo"
)

const (
	ddstatsChannelName = "ddstats"
	prefix             = "."
)

type Discord struct {
	Session         *discordgo.Session
	DB              *postgres.Postgres
	ddAPI           *ddapi.API
	websocketHub    *websocket.Hub
	commands        *sync.Map
	ddstatsChannels *ddstatsChannels
	infoLog         *log.Logger
	errorLog        *log.Logger
}

func New(token string, db *postgres.Postgres, ddAPI *ddapi.API, websocketHub *websocket.Hub, infoLog, errorLog *log.Logger) (*Discord, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	discord := Discord{
		Session:         session,
		DB:              db,
		ddAPI:           ddAPI,
		websocketHub:    websocketHub,
		commands:        &sync.Map{},
		ddstatsChannels: &ddstatsChannels{},
		infoLog:         infoLog,
		errorLog:        errorLog,
	}
	session.AddHandler(discord.messageCreate)
	discord.registerCommands()
	return &discord, nil
}

func (d *Discord) Start() error {
	d.infoLog.Println("Starting Discord Bot")
	err := d.Session.Open()
	if err != nil {
		return err
	}
	err = d.getDDStatsChannels()
	if err != nil {
		return err
	}
	err = d.Session.UpdateStatusComplex(discordgo.UpdateStatusData{Game: &discordgo.Game{Name: ".help | ddstats.com"}})
	if err != nil {
		return err
	}
	return nil
}

func (d *Discord) Close() {
	d.Session.Close()
}

func (d *Discord) getDDStatsChannels() error {
	for _, guild := range d.Session.State.Guilds {
		channel, err := d.Session.GuildChannels(guild.ID)
		if err != nil {
			return err
		}
		for _, c := range channel {
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			if strings.Contains(c.Name, ddstatsChannelName) {
				d.ddstatsChannels.store(c.ID)
			}
		}

	}
	return nil
}

type ddstatsChannels struct {
	sync.Mutex
	channels []string
}

func (ddc *ddstatsChannels) store(id string) {
	ddc.Lock()
	defer ddc.Unlock()
	ddc.channels = append(ddc.channels, id)
}

func (ddc *ddstatsChannels) load() []string {
	ddc.Lock()
	defer ddc.Unlock()
	return ddc.channels
}

func (ddc *ddstatsChannels) contains(id string) bool {
	ddc.Lock()
	defer ddc.Unlock()
	for _, c := range ddc.channels {
		if c == id {
			return true
		}
	}
	return false
}
