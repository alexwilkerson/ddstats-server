package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexwilkerson/ddstats-server/pkg/collector"
	"github.com/alexwilkerson/ddstats-server/pkg/ddapi"

	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
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

	postgresDB := postgres.NewPostgres(client, db)

	ddAPI := ddapi.NewAPI(client)

	collector := collector.NewCollector(ddAPI, postgresDB, infoLog, errorLog)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	infoLog.Println("Starting the collector...")
	start := time.Now()
	go collector.Start()
	go func() {
		<-quit
		fmt.Println("\r") // overwrites ^C char
		infoLog.Println("Shutting down the collector...")
		collector.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-collector.Done():
	}
	fmt.Println(time.Since(start))
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
