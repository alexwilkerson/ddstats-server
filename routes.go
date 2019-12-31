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

	return standardMiddleware.Then(mux)
}
