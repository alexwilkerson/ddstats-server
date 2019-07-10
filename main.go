package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexwilkerson/ddstats-api/server"
)

var (
	CertFile      = os.Getenv("CERT_FILE")
	KeyFile       = os.Getenv("KEY_FILE")
	ServerAddress = os.Getenv("SERVER_ADDRESS")
)

const (
	message = "Hello, World."
)

func main() {
	// logger := log.New(os.Stdout, "[ddstats-api] ", log.LstdFlags)
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)

	err := server.New(mux, ServerAddress).ListenAndServeTLS(CertFile, KeyFile)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
