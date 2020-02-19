package discord

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"

	"github.com/alexwilkerson/ddstats-server/pkg/websocket"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

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
	quit            chan struct{}
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
		quit:            make(chan struct{}),
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
	go d.listenForNotifications()
	return nil
}

func (d *Discord) listenForNotifications() {
	for {
		select {
		case notification := <-d.websocketHub.DiscordBroadcast:
			switch v := notification.(type) {
			case *websocket.PlayerBestReached:
				go func() {
					err := d.broadcast(&discordgo.MessageEmbed{
						Title:       fmt.Sprintf("%s just passed their best time of %.4fs!", v.PlayerName, v.PreviousGameTime),
						Description: fmt.Sprintf("Watch here: https://ddstats.com/players/%d", v.PlayerID),
					})
					if err != nil {
						d.errorLog.Printf("%+v", err)
					}
				}()
			case *websocket.PlayerBestSubmitted:
				go func() {
					err := d.broadcast(&discordgo.MessageEmbed{
						Title:       fmt.Sprintf("%s just got a new score of %.4fs!", v.PlayerName, v.GameTime),
						Description: fmt.Sprintf("...beating their old high score of %.4fs by %.4f seconds!\nGame log here: https://ddstats.com/games/%d", v.PreviousGameTime, v.GameTime-v.PreviousGameTime, v.GameID),
					})
					if err != nil {
						d.errorLog.Printf("%+v", err)
					}
				}()
			case *websocket.PlayerAboveThreshold:
				go func() {
					err := d.broadcast(&discordgo.MessageEmbed{
						Title:       fmt.Sprintf("%s is above 1000!", v.PlayerName),
						Description: fmt.Sprintf("Watch here: https://ddstats.com/players/%d", v.PlayerID),
					})
					if err != nil {
						d.errorLog.Printf("%+v", err)
					}
				}()
			case *websocket.PlayerDied:
				go func() {
					err := d.broadcast(&discordgo.MessageEmbed{
						Title:       fmt.Sprintf("%s died at %.4f", v.PlayerName, v.GameTime),
						Description: fmt.Sprintf("...%s\nGame log: https://ddstats.com/games/%d", strings.ToLower(v.DeathType), v.GameID),
					})
					if err != nil {
						d.errorLog.Printf("%+v", err)
					}
				}()
			default:
				d.errorLog.Println("invalid type received to discord listener")
			}
		case <-d.quit:
			return
		}
	}
}

func (d *Discord) broadcast(embed *discordgo.MessageEmbed) error {
	embed.Color = defaultColor
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    "ddstats.com",
		IconURL: iconURL,
	}
	for _, channel := range d.ddstatsChannels.load() {
		_, err := d.Session.ChannelMessageSendEmbed(channel, embed)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Discord) Close() {
	close(d.quit)
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
