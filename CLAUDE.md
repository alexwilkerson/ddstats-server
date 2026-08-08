# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Backend server for [ddstats.com](https://ddstats.com), a stats-tracking site for the game Devil Daggers. It ingests game runs submitted by a companion desktop client, proxies/caches data from the official Devil Daggers backend ("DD API"), serves a REST API + WebSocket feed to a Vue frontend (`ui/`), runs a Discord bot for live notifications, and periodically crawls the DD leaderboard to compute aggregate stats.

Two binaries share the same `pkg/` code:
- `cmd/server` — the long-running API/websocket/socket.io/gRPC/Discord server
- `cmd/collector` — a one-shot process (run via cron) that crawls the entire DD leaderboard and records a `collector_run` snapshot

## Commands

```bash
# Run the server locally (uses reflex to auto-restart on .go changes)
reflex -d none -c reflex.conf
# or directly:
go run ./cmd/server --dsn "host=localhost port=5432 user=ddstats password=ddstats dbname=ddstats sslmode=disable"

# Run the collector once
go run ./cmd/collector --dsn "..."

# Build both binaries for linux deploy target
make build

# Deploy (scp to remote host "casd") — do not run without confirming with the user
make deploy

# Tests (only pkg/api and pkg/discord have tests currently)
go test ./...
go test ./pkg/api/...
go test -run TestName ./pkg/discord/...

# Regenerate gRPC code after editing gamesubmission/gamesubmission.proto
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  gamesubmission/gamesubmission.proto
```

Frontend (`ui/`, Vue 2 + Vuetify) has its own toolchain:
```bash
cd ui && yarn serve   # dev server
cd ui && yarn build   # production build, served by the Go server from ui/dist
cd ui && yarn lint
```

`.github/workflows/ci.yml` runs `gofmt`/`go vet`/`go build`/`go test -race` and the UI lint+build on every PR/push to `master`; both checks are required on `master` via branch protection. `.github/workflows/release.yml` builds and deploys `cmd/server`, `cmd/collector`, and `ui/dist` to the production host when a GitHub Release is published (or via manual `workflow_dispatch`).

## Architecture

### Request flow / transport multiplexing

`cmd/server/main.go` opens a single TCP listener and uses `cmux` to demux three protocols on one port:
1. gRPC (matched by the `application/grpc` content-type header) — implements `gamesubmission.GameRecorder` (see `cmd/server/grpc.go`), used by the newer desktop client to submit games and check for client updates/MOTD.
2. Everything else goes to a stdlib `http.Server` whose handler is built by `pkg/api.Routes()`.

Inside the HTTP handler, routing is itself split in two:
- A parent `http.ServeMux` handles `/socket.io/` (legacy live-stats protocol, backward-compat for older clients) and serves the built Vue app as static files for every other path.
- `/api/` and `/api/v2/` are delegated to a `pat.Mux` (`github.com/bmizerany/pat`) wrapped in `alice` middleware (`recoverPanic`, `handleCORS`, `logRequest`, `secureHeaders`) — see `pkg/api/routes.go`.
- `/ws` is a raw `gorilla/websocket` upgrade handled by `pkg/websocket`.

When changing routing, note the comment in `routes.go` explaining *why* the split exists: `pat` only handles REST-style paths, so non-REST traffic (websockets, the SPA) has to be handled by the outer mux first.

### Two live-update mechanisms (legacy + current)

There are two independent real-time systems that both exist because of backward compatibility with older desktop clients:
- `pkg/socketio` — the older `go-socket.io`-based protocol (events like `login`, `submit`, `state_update`, `game_submitted`). New clients (0.6.1+) use `state_update` with `leviDownTime`/`orbDownTime`; `onSubmit` is a compatibility shim for pre-0.6.1 clients that calls into `onStateUpdate` with those defaulted to 0.
- `pkg/websocket` — a hand-rolled hub/client/room model (`Hub.Start()` is a single goroutine handling all events via channels: `Register`, `Broadcast`, `RegisterPlayer`, `DiscordBroadcast`, etc.) that pushes JSON messages to browser clients on `/ws` and to the Discord bot.

`pkg/socketio` is the producer for most live events; it pushes into `pkg/websocket`'s `Hub` channels (`Broadcast`, `BroadcastToAll`, `DiscordBroadcast`) rather than talking to browsers directly. When adding a new live event, decide whether it needs to go out over socket.io (legacy clients), the websocket hub (browser + Discord), or both.

### Data layer

`pkg/models` defines plain structs (DB row shapes and API request/response shapes — note `SubmittedGame` vs `SubmittedGameV2`, the JSON shapes for old vs new client submission formats). `pkg/models/postgres` implements one `*Model` type per table (e.g. `GameModel`, `PlayerModel`, `CollectorRunModel`) wrapping `*sqlx.DB`/`*sqlx.Tx`, aggregated into a single `Postgres` struct (`pkg/models/postgres/postgres.go`) that's constructed once in `main` and passed everywhere by pointer. There's no ORM — write raw SQL against `schema.sql`.

`schema.sql` is destructive-first: it starts with `DROP TABLE` statements before `CREATE TABLE IF NOT EXISTS`, so it's meant to be run against a fresh/disposable database, not applied as a migration to production.

### DD API integration (`pkg/ddapi`)

The upstream Devil Daggers backend (`dd.hasmodai.com`) has no JSON API — it returns raw binary blobs that `pkg/ddapi` parses by hand with hardcoded byte offsets (`bytesToPlayer`, `bytesToLeaderboard`, `userSearchBytesToPlayers`). If DD ever changes their response format, these offsets need to be re-derived; there's no schema to consult.

### Collector (`pkg/collector`)

`Collector.Start()` pages through the *entire* DD leaderboard (`maxLimit = 1000` per page) inside a single DB transaction, diffing each player against their previously stored `CollectorPlayer` row to compute deltas (new deaths, rank/score improvement, dagger-tier crossings at 60/120/250/500s thresholds) and rolling them up into one `CollectorRun` row. It's designed to be safely interruptible (`Stop()`/`quit` channel triggers a `tx.Rollback()`), since a full crawl can take a while and is invoked externally (e.g. cron), not by the server process.

### Discord bot (`pkg/discord`)

Commands live one-per-file (`command_*.go`) and self-register into a `sync.Map` via `registerCommands()`/`command.go`; `discord.go` is the client-agnostic subscriber to `websocketHub.DiscordBroadcast` that turns hub notification events (`PlayerBestReached`, `PlayerAboveThreshold`, etc.) into embeds posted to any channel whose name contains `ddstats` in every guild the bot is in. Adding a new notification type means adding a struct + channel send in `pkg/websocket`, a case in `discord.go`'s `listenForNotifications` switch, and usually a corresponding push out of `pkg/socketio`.

### gRPC / protobuf (`gamesubmission/`)

`gamesubmission.proto` is the source of truth for the newer submission protocol; `.pb.go` / `_grpc.pb.go` are generated (see Commands above) — edit the `.proto` and regenerate rather than hand-editing the generated files.
