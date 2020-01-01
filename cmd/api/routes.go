package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.handleCors, app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.helloWorld))

	// ddapi
	mux.Get("/api/v2/ddapi/get_user_by_rank", http.HandlerFunc(app.ddGetUserByRank))
	mux.Get("/api/v2/ddapi/get_user_by_id", http.HandlerFunc(app.ddGetUserByID))
	mux.Get("/api/v2/ddapi/get_user_by_name", http.HandlerFunc(app.ddUserSearch))
	mux.Get("/api/v2/ddapi/get_scores", http.HandlerFunc(app.ddGetScores))
	mux.Get("/api/v2/game", http.HandlerFunc(app.getGame))

	return standardMiddleware.Then(mux)
}
