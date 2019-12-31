package main

import (
	"net/http"
	"strconv"
)

func (app *application) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Merle!"))
}

func (app *application) ddGetUserByRank(w http.ResponseWriter, r *http.Request) {
	rank, ok := r.URL.Query()["rank"]
	if !ok || len(rank) < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	rankInt, err := strconv.Atoi(rank[0])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if rankInt < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	w.Write([]byte(strconv.Itoa(rankInt)))

}
