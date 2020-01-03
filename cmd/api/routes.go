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

	// ddstats api
	// this endpoint is redundant so as to handle older client submissions
	mux.Post("/api/submit_game", http.HandlerFunc(app.submitGame))
	mux.Post("/api/v2/submit_game", http.HandlerFunc(app.submitGame))
	mux.Get("/api/v2/game/top", http.HandlerFunc(app.getTopGames))
	mux.Get("/api/v2/game/recent", http.HandlerFunc(app.getRecentGames))
	mux.Get("/api/v2/game", http.HandlerFunc(app.getGame))
	mux.Get("/api/v2/game/all", http.HandlerFunc(app.getGameAll))
	mux.Get("/api/v2/game/gems", http.HandlerFunc(app.getGameGems))
	mux.Get("/api/v2/game/homing-daggers", http.HandlerFunc(app.getGameHomingDaggers))
	mux.Get("/api/v2/game/daggers-hit", http.HandlerFunc(app.getGameDaggersHit))
	mux.Get("/api/v2/game/daggers-fired", http.HandlerFunc(app.getGameDaggersFired))
	mux.Get("/api/v2/game/accuracy", http.HandlerFunc(app.getGameAccuracy))
	mux.Get("/api/v2/game/enemies-alive", http.HandlerFunc(app.getGameEnemiesAlive))
	mux.Get("/api/v2/game/enemies-killed", http.HandlerFunc(app.getGameEnemiesKilled))
	mux.Get("/api/v2/player", http.HandlerFunc(app.getPlayer))
	mux.Get("/api/v2/player/all", http.HandlerFunc(app.getPlayers))

	return standardMiddleware.Then(mux)
}
