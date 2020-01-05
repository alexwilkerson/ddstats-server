package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexwilkerson/ddstats-api/pkg/ddapi"
	"github.com/alexwilkerson/ddstats-api/pkg/models/postgres"
	socketio "github.com/googollee/go-socket.io"
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
	ddAPI          *ddapi.API
	games          *postgres.GameModel
	players        *postgres.PlayerModel
	submittedGames *postgres.SubmittedGameModel
	motd           *postgres.MOTDModel
	socketIO       *socketio.Server
}

func main() {
	addr := flag.String("addr", ":5000", "HTTP Network Address")
	dsn := flag.String("dsn", "host=localhost port=5432 user=ddstats password=ddstats dbname=ddstats sslmode=disable", "PostgreSQL data source name")
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

	socketIOServer, err := socketio.NewServer(nil)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		client:         client,
		ddAPI:          &ddapi.API{Client: client},
		games:          &postgres.GameModel{DB: db},
		players:        &postgres.PlayerModel{DB: db},
		submittedGames: &postgres.SubmittedGameModel{DB: db, Client: client},
		motd:           &postgres.MOTDModel{DB: db},
		socketIO:       socketIOServer,
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
