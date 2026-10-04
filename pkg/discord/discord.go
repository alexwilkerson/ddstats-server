package discord

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"

	"github.com/alexwilkerson/ddstats-server/pkg/websocket"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	"github.com/bwmarrin/discordgo"
	gorillawebsocket "github.com/gorilla/websocket"
)

const (
	ddstatsChannelName = "ddstats"
	prefix             = "."
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes
	closeCodeDisallowedIntents = 4014
)

type Discord struct {
	Session         *discordgo.Session
	DB              *postgres.Postgres
	ddAPI           *ddapi.API
	websocketHub    *websocket.Hub
	commands        *sync.Map
	ddstatsChannels *ddstatsChannels
	ddStatus        *ddStatus
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
		ddStatus:        &ddStatus{},
		infoLog:         infoLog,
		errorLog:        errorLog,
		quit:            make(chan struct{}),
	}
	// Message Content is a privileged intent: it must also be enabled for the
	// bot in the Discord Developer Portal, or the gateway rejects the
	// connection with close code 4014.
	session.Identify.Intents = discordgo.IntentGuilds |
		discordgo.IntentGuildMessages |
		discordgo.IntentDirectMessages |
		discordgo.IntentMessageContent
	session.LogLevel = discordgo.LogWarning
	discordgo.Logger = discordgoLogger(infoLog, errorLog)
	session.AddHandler(discord.messageCreate)
	session.AddHandler(discord.onReady)
	session.AddHandler(discord.onResumed)
	session.AddHandler(discord.onDisconnect)
	discord.registerCommands()
	return &discord, nil
}

func (d *Discord) Start() error {
	d.infoLog.Println("Starting Discord Bot")
	err := d.Session.Open()
	var closeErr *gorillawebsocket.CloseError
	if errors.As(err, &closeErr) && closeErr.Code == closeCodeDisallowedIntents {
		return fmt.Errorf("opening discord gateway connection: %w (enable the Message Content intent for the bot in the Discord Developer Portal)", err)
	}
	if err != nil {
		return fmt.Errorf("opening discord gateway connection: %w", err)
	}
	err = d.getDDStatsChannels()
	if err != nil {
		return fmt.Errorf("finding ddstats channels: %w", err)
	}
	d.infoLog.Printf("Discord bot broadcasting to %d ddstats channel(s)", len(d.ddstatsChannels.load()))
	err = d.Session.UpdateGameStatus(0, ".help | ddstats.com")
	if err != nil {
		return fmt.Errorf("setting discord status: %w", err)
	}
	go d.listenForNotifications()
	go d.monitorDDStatus()
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
			case *websocket.PlayerAboveThresholdSubmitted:
				go func() {
					err := d.broadcast(&discordgo.MessageEmbed{
						Title:       fmt.Sprintf("%s died at %.4f seconds!", v.PlayerName, v.GameTime),
						Description: fmt.Sprintf("...%s\nGame log here: https://ddstats.com/games/%d", strings.ToLower(v.DeathType), v.GameID),
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
			return fmt.Errorf("broadcasting %q to channel %s: %w", embed.Title, channel, err)
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
