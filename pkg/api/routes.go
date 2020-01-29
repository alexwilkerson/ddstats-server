package api

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"

	socketio "github.com/googollee/go-socket.io"
)

func (api *API) Routes(socketioServer *socketio.Server) http.Handler {
	standardMiddleware := alice.New(api.recoverPanic, api.handleCORS, api.logRequest, secureHeaders)

	mux := pat.New()

	// ddapi
	mux.Get("/api/v2/ddapi/get_user_by_rank", http.HandlerFunc(api.ddGetUserByRank))
	mux.Get("/api/v2/ddapi/get_user_by_id", http.HandlerFunc(api.ddGetUserByID))
	mux.Get("/api/v2/ddapi/get_user_by_name", http.HandlerFunc(api.ddUserSearch))
	mux.Get("/api/v2/ddapi/get_scores", http.HandlerFunc(api.ddGetScores))

	// ddstats api
	mux.Post("/api/v2/submit_game", http.HandlerFunc(api.submitGame))
	mux.Post("/api/v2/client_connect", http.HandlerFunc(api.clientConnect))
	mux.Get("/api/v2/game/top", http.HandlerFunc(api.getTopGames))
	mux.Get("/api/v2/leaderboard", http.HandlerFunc(api.getLeaderboard))
	mux.Get("/api/v2/game/recent", http.HandlerFunc(api.getRecentGames))
	mux.Get("/api/v2/game", http.HandlerFunc(api.getGame))
	mux.Get("/api/v2/game/full", http.HandlerFunc(api.getGameFull))
	mux.Get("/api/v2/game/all", http.HandlerFunc(api.getGameAll))
	mux.Get("/api/v2/game/gems", http.HandlerFunc(api.getGameGems))
	mux.Get("/api/v2/game/homing_daggers", http.HandlerFunc(api.getGameHomingDaggers))
	mux.Get("/api/v2/game/daggers_hit", http.HandlerFunc(api.getGameDaggersHit))
	mux.Get("/api/v2/game/daggers_fired", http.HandlerFunc(api.getGameDaggersFired))
	mux.Get("/api/v2/game/accuracy", http.HandlerFunc(api.getGameAccuracy))
	mux.Get("/api/v2/game/enemies_alive", http.HandlerFunc(api.getGameEnemiesAlive))
	mux.Get("/api/v2/game/enemies_killed", http.HandlerFunc(api.getGameEnemiesKilled))
	mux.Get("/api/v2/player", http.HandlerFunc(api.getPlayer))
	mux.Get("/api/v2/player/update", http.HandlerFunc(api.playerUpdate))
	mux.Get("/api/v2/player/live", http.HandlerFunc(api.playerLive))
	mux.Get("/api/v2/player/all", http.HandlerFunc(api.getPlayers))
	mux.Get("/api/v2/motd", http.HandlerFunc(api.getMOTD))
	mux.Get("/api/v2/releases", http.HandlerFunc(api.getReleases))
	mux.Get("/api/v2/news", http.HandlerFunc(api.getNews))
	mux.Get("/api/v2/daily", http.HandlerFunc(api.getDaily))

	// these are here for now to be backward compatible
	mux.Post("/api/get_motd", http.HandlerFunc(api.clientConnect))
	mux.Post("/api/submit_game", http.HandlerFunc(api.submitGame))

	mux.Get("/ws", api.handleCORS(http.HandlerFunc(api.serveWebsocket)))

	// Why? Well, because the pat application only accounts for REST requests,
	// so if the server receives anything else (such as a websocket request),
	// there's no way to register it.. these three lines will match the /socket-io/
	// end point and if it doesn't match will pass everything on to the pat mux
	// since "/" matches everything
	muxParent := http.NewServeMux()
	muxParent.Handle("/socket.io/", socketioCORS(socketioServer))

	// serves the vue app, built inside the dist directory
	// must be handled by this parent mux, since for whatever
	// reason it won't work otherwise
	vueApp := api.handleCORS(http.FileServer(http.Dir("./ui/dist/")))
	muxParent.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(r.URL.RequestURI(), vueApp).ServeHTTP(w, r)
	}))
	// these routes are needed to point to specific static files generated
	// by vue
	muxParent.Handle("/js/", http.StripPrefix("/", vueApp))
	muxParent.Handle("/fonts/", http.StripPrefix("/", vueApp))
	muxParent.Handle("/css/", http.StripPrefix("/", vueApp))
	muxParent.Handle("/img/", http.StripPrefix("/", vueApp))
	muxParent.Handle("/static/", http.StripPrefix("/", vueApp))
	muxParent.Handle("/favicon.ico", http.StripPrefix("/", vueApp))
	muxParent.Handle("/android-chrome-192x192.png", http.StripPrefix("/", vueApp))
	muxParent.Handle("/android-chrome-512x512.png", http.StripPrefix("/", vueApp))
	muxParent.Handle("/apple-touch-icon.png", http.StripPrefix("/", vueApp))
	muxParent.Handle("/favicon-16x16.png", http.StripPrefix("/", vueApp))
	muxParent.Handle("/favicon-32x32.png", http.StripPrefix("/", vueApp))
	muxParent.Handle("/site.webmanifest", http.StripPrefix("/", vueApp))
	// END VUEJS BULLSHIT

	muxParent.Handle("/api/", standardMiddleware.Then(mux))
	muxParent.Handle("/api/v2/", standardMiddleware.Then(mux))
	muxParent.Handle("/ws/", standardMiddleware.Then(mux))

	return muxParent
}
