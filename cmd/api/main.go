package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/discord"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
	"github.com/alexwilkerson/ddstats-api/pkg/models/postgres"
	"github.com/alexwilkerson/ddstats-api/pkg/socketio"
	"github.com/alexwilkerson/ddstats-api/pkg/websocket"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	oldestValidClientVersion = "0.3.1"
	currentClientVersion     = "0.4.5"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	client         *http.Client
	websocketHub   *websocket.Hub
	ddAPI          *ddapi.API
	games          *postgres.GameModel
	players        *postgres.PlayerModel
	submittedGames *postgres.SubmittedGameModel
	motd           *postgres.MOTDModel
}

func main() {
	addr := flag.String("addr", ":5000", "HTTP Network Address")
	dsn := flag.String("dsn", "host=localhost port=5432 user=ddstats password=ddstats dbname=ddstats sslmode=disable", "PostgreSQL data source name")
	discordToken := flag.String("discord-token", "NjY1MDY4MDcwOTQ3MTI3MzA3.XhgWNQ.W60yL9JnPbKNFUSz1XEbNpuoYs8", "Discord Bot Token")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// TODO: set up client appropriately
	client := &http.Client{}

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		client:         client,
		websocketHub:   websocket.NewHub(),
		ddAPI:          &ddapi.API{Client: client},
		games:          &postgres.GameModel{DB: db},
		players:        &postgres.PlayerModel{DB: db},
		submittedGames: &postgres.SubmittedGameModel{DB: db, Client: client},
		motd:           &postgres.MOTDModel{DB: db},
	}

	socketioServer, err := socketio.NewServer(infoLog, errorLog, app.websocketHub, client, db)
	if err != nil {
		errorLog.Fatal(err)
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(socketioServer),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	discordSession, err := discord.New(*discordToken, infoLog, errorLog)
	if err != nil {
		errorLog.Fatal(err)
	}
	err = discordSession.Start()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer discordSession.Close()
	go app.websocketHub.Start()
	go socketioServer.Serve()
	defer socketioServer.Close()

	infoLog.Printf("Starting server on %s", *addr)
	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
