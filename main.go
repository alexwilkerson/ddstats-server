package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexwilkerson/ddstats-api/server"
	"github.com/gorilla/mux"
)

var (
	CertFile      = os.Getenv("CERT_FILE")
	KeyFile       = os.Getenv("KEY_FILE")
	ServerAddress = os.Getenv("SERVER_ADDRESS")
)

const (
	message = "Hello, World."
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	staticFileDirectory := http.Dir("./static")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")
	r.HandleFunc("/", HomeHandler).Methods("GET")
	return r
}

func main() {
	// logger := log.New(os.Stdout, "[ddstats-api] ", log.LstdFlags)
	r := newRouter()

	err := server.New(r, ServerAddress).ListenAndServeTLS(CertFile, KeyFile)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
