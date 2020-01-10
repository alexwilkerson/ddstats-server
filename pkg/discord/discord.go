package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Session  *discordgo.Session
	infoLog  *log.Logger
	errorLog *log.Logger
}

func New(token string, infoLog, errorLog *log.Logger) (*Discord, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	discord := Discord{
		Session:  session,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
	discord.Session.AddHandler(discord.messageCreate)
	return &discord, nil
}

func (d *Discord) Start() error {
	d.infoLog.Println("Starting Discord Bot")
	err := d.Session.Open()
	if err != nil {
		return err
	}
	err = d.Session.UpdateListeningStatus("666")
	if err != nil {
		return err
	}
	return nil
}

func (d *Discord) Close() {
	d.Session.Close()
}
