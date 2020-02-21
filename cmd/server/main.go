package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"

	"github.com/alexwilkerson/ddstats-server/pkg/api"
	"github.com/alexwilkerson/ddstats-server/pkg/discord"

	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"
	"github.com/alexwilkerson/ddstats-server/pkg/socketio"
	"github.com/alexwilkerson/ddstats-server/pkg/websocket"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	client       *http.Client
	websocketHub *websocket.Hub
	ddAPI        *ddapi.API
}

func main() {
	addr := flag.String("addr", ":5000", "HTTP Network Address")
	dsn := flag.String("dsn", "host=localhost port=5432 user=ddstats password=ddstats dbname=ddstats sslmode=disable", "PostgreSQL data source name")
	discordToken := flag.String("discord-token", "***REMOVED***", "Discord Bot Token")
	disableDiscord := flag.Bool("disable-discord", false, "Disable the Discord Bot")
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

	postgresDB := postgres.NewPostgres(client, db)

	websocketHub := websocket.NewHub(postgresDB)

	ddAPI := ddapi.NewAPI(client)

	api := api.NewAPI(client, postgresDB, websocketHub, ddAPI, infoLog, errorLog)

	socketioServer, err := socketio.NewServer(infoLog, errorLog, websocketHub, client, postgresDB)
	if err != nil {
		errorLog.Fatal(err)
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      api.Routes(socketioServer),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if !*disableDiscord {
		discordSession, err := discord.New(*discordToken, postgresDB, ddAPI, websocketHub, infoLog, errorLog)
		if err != nil {
			errorLog.Fatal(err)
		}
		err = discordSession.Start()
		if err != nil {
			errorLog.Fatal(err)
		}
		defer discordSession.Close()
	}

	go websocketHub.Start()
	defer websocketHub.Close()
	go socketioServer.Serve()
	defer socketioServer.Close()

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		infoLog.Println("Server shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			errorLog.Fatal("Could not gracefully shut down the server: %w", err)
		}
		close(done)
	}()

	infoLog.Printf("Starting server on %s", *addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errorLog.Fatalf("could not listen on %s: %w", *addr, err)
	}

	<-done
	infoLog.Println("Server stopped")
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
