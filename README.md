# ddstats-server

![DDSTATS Server v2](/server_diagram.png)

## Summary

Date : 2020-01-12 21:02:43

Total : 42 files, 3887 codes, 176 comments, 528 blanks, all 4591 lines

### Languages

| language   | files |  code | comment | blank | total |
| :--------- | ----: | ----: | ------: | ----: | ----: |
| Go         |    37 | 3,597 |     151 |   464 | 4,212 |
| Markdown   |     1 |   139 |       0 |    45 |   184 |
| SQL        |     1 |   109 |       3 |    14 |   126 |
| HTTP       |     1 |    25 |      21 |     2 |    48 |
| XML        |     1 |    16 |       0 |     3 |    19 |
| Properties |     1 |     1 |       1 |     0 |     2 |

### Directories

| path                | files |  code | comment | blank | total |
| :------------------ | ----: | ----: | ------: | ----: | ----: |
| .                   |    42 | 3,887 |     176 |   528 | 4,591 |
| cmd                 |     1 |   101 |       1 |    26 |   128 |
| cmd/server          |     1 |   101 |       1 |    26 |   128 |
| pkg                 |    36 | 3,496 |     150 |   438 | 4,084 |
| pkg/api             |     7 |   817 |      16 |   155 |   988 |
| pkg/ddapi           |     2 |   303 |      34 |    54 |   391 |
| pkg/discord         |    14 |   888 |      15 |    88 |   991 |
| pkg/models          |     8 |   915 |      57 |    84 | 1,056 |
| pkg/models/postgres |     7 |   778 |      41 |    66 |   885 |
| pkg/socketio        |     1 |   266 |       9 |    25 |   300 |
| pkg/websocket       |     4 |   307 |      19 |    32 |   358 |

## Building for linux

```
GOOS=linux GOARCH=amd64 go build -o dist/server -v ./cmd/server
```

```
GOOS=linux GOARCH=amd64 go build -o dist/collector -v ./cmd/collector
```

## Automatically restarting server during dev

Go 1.11+ installation required

```
go get github.com/cespare/reflex
```

Run reflex using the following command:

```
reflex -d none -c reflex.conf
```

## Things Left to Do...

- [ ] authorization

#### ui

##### pages

- [ ] index
  - collector statistics
  - global stats
  - live users
  - news
- [ ] players
  - live users
- [ ] info
  - live users
- [ ] individual player page
  - live users
  - live stats
  - recent games
  - player global stats
- [ ] download page
  - live users
- [ ] leaderboard page
  - filter by leaderboard
  - live users
- [ ] game stats
  - live users

#### collector

- [x] add game_time_improvement to collector_active_player
- [x] add the following to runs: bronze_daggers, silver_daggers, gold_daggers, devil_daggers, total_bronze, total_silver, total_gold, total_devil
- [x] add last_active to player table and collect using collector
- [x] possibly add deaths and total_game_time to active players column

#### discord-bot commands

_MAKE SURE TO CHECK FOR SURVIVAL HASH WHEN MAKING THE NOTIFICATION FUNCTION!_

- [x] global https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/global.js
- [x] help https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/help.js
- [x] id https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/id.js
- [x] live https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/live.js
- [x] rank https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/rank.js
- [x] register https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/register.js
- [x] me https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/me.js
- [x] search https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/search.js
- [x] top https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/top.js

#### submission stuff

- [x] finish the logic behind game submissions and filtering (i think this is done)

#### live notification functionality

- [x] verify & fix bugs in client that relate to game submissions https://github.com/alexwilkerson/ddstats-go/blob/master/net.go#L80-L86
- [x] create socket.io listener for `game_submitted` function which has arguments (`gameID int`, `notifyPlayerBest bool`, `notifyAbove1000 bool`)
- [x] create discord bot listener which runs in its own goroutine and has channels for notifying player best and notifying when player is above 1000
- [x] create connection between socket.io -> websocket -> discord for player best and above 1000 notifications

### TODO: endpoints

#### collector

- [x] `api/v2/daily` returns most recent
- [ ] `api/v2/daily?date={?}`

#### index

- [x] set up `/static` endpoint to host static files
- [x] GET player/live (gets live players from websocket)
- [x] POST api/v2/client_connect (client sends client version and receives info)
- [x] GET api/v2/motd
- [x] create table in database for motd
- [x] GET api/v2/client/releases?page_size={int}&page_num={int}
- [x] create table for client releases
- [x] GET api/v2/news?page_size={int}&page_num={int}
- [x] create table in database for news

#### game

- [x] POST api/v2/game/submit
- [x] GET api/v2/game/top
- [x] GET api/v2/leaderboard?spawnset={spawnset_name}?page_size={int}&page_num={int}
- [x] GET api/v2/game/recent?page_size={int}&page_num={int} (optional: player_id={int} will give most recent for given player)
- [x] GET api/v2/game?id={int} (info)
- [x] GET api/v2/game/all?id={int}
- [x] GET api/v2/game/gems?id={int}
- [x] GET api/v2/game/homing-daggers?id={int}
- [x] GET api/v2/game/accuracy?id={int}
- [x] GET api/v2/game/enemies-alive?id={int}
- [x] GET api/v2/game/enemies-killed?id={int}

#### player

- [x] GET api/v2/player/update?id={int} (updates or creates user in database from dd backend)
- [x] GET api/v2/player/all?page_size={int}&page_num={int} (list all users)
- [x] GET api/v2/player?id={int}
- [x] ~GET api/v2/player/games?id={int}page_size={int}&page_num={int} (paginated list of games by player id)~
      (see api/v2/game/recent)

#### dd backend

- [x] GET api/v2/ddapi/get_user_by_rank?rank={int}
- [x] GET api/v2/ddapi/get_user_by_id?id={int}
- [x] GET api/v2/ddapi/user_search?user={string}
- [x] GET api/v2/ddapi/get_scores?offset={int}&limit={int}

## Experimental stuff

- store daily data about ddstats users
  - each day, download their information to the database
  - end goal would be to eventually show stats over a period of time so users could continuously track their improvement day-to-day
  - possible daily data:
    - all data for that day
      - play time
      - number of deaths
      - gems
      - daggers fired
      - daggers hit
      - enemies killed
      - accuracy
      - etc, etc, etc... basically all of the player data
- a function which runs every day at midnight which goes through the entire devil daggers backend, retrieves user information, stores it to the database. during the process of data collection, the function will analyze the informatiod and store the data to the database. it's possible that this function will be too expensive to process or take too much time, but hopefully it would work.
  - the function would run daily
  - possible daily data:
    - users active for that day
    - how many players got a new score
    - how many new players
    - average improvement among players who got new scores

## Exporting data from SQLite3 database to PostgreSQL database

### importing csvs from app.db

In order to port the database from SQLite3 to Postgres, the following shenanigans must occur:

- scp app.db from server, then: `sqlite3 app.db`
- Afterwards, import the csvs via Postico.
- See `schema.sql` for creating the Postgres database.

#### game

```sql
.headers on
.mode csv
.output game.csv
-- printf required because of sqlite's weird float formatting
select id, player_id, granularity, printf("%.6f", game_time) as game_time, death_type, gems, homing_daggers, daggers_fired, daggers_hit, enemies_alive, enemies_killed, time_stamp, replay_player_id, survival_hash, version, printf("%.6f", level_two_time) as level_two_time, printf("%.6f", level_three_time) as level_three_time, printf("%.6f", level_four_time) as level_four_time, printf("%.6f", homing_daggers_max_time) as homing_daggers_max_time, printf("%.6f", enemies_alive_max_time) as enemies_alive_max_time, homing_daggers_max, enemies_alive_max from game;
```

#### state

```sql
.headers on
.mode csv
.output state.csv
-- remove any inconsistencies between state and game tables
delete from state where game_id not in (select id from game where id is not null);
-- printf required because of sqlite's weird float formatting
select id, game_id, printf("%.6f", game_time) as game_time, gems, homing_daggers, daggers_hit, daggers_fired, enemies_alive, enemies_killed from state;
```

might have to split up state to get it into postgres
`split -l 2500000 state.csv` works

#### spawnset

```sql
.headers on
.mode csv
.output spawnset.csv
select * from spawnset;
```

#### player (renamed from user)

```sql
.headers on
.mode csv
.output player.csv
select id, username as player_name, rank, game_time, death_type, gems, daggers_fired, daggers_hit, enemies_killed, accuracy, time_total as overall_game_time, deaths_total as overall_deaths, gems_total as overall_gems, enemies_killed_total as overall_enemies_killed, daggers_fired_total as overall_daggers_fired, daggers_hit_total as overall_daggers_hit, accuracy_total as overall_accuracy from user;
```

add overall_game_time: `UPDATE player SET overall_average_game_time=TRUNC(DIVZERO(overall_game_time, overall_deaths)::numeric, 4);`

#### live

```sql
.headers on
.mode csv
.output live.csv
select * from live;
```

### discord bot database

#### discord_users

```sql
.headers on
.mode csv
.output discord_users.csv
select * from users;
```

### little helper to delete non-ddstats-users from database

````sql
DELETE FROM player
WHERE id<>-1 AND id IN (
	SELECT player.id
	FROM player LEFT JOIN game ON player.id=game.player_id
	WHERE game.id IS NULL
	GROUP BY player.id
);```

![DDSTATS Logo](/ui/static/ddstats_logo_v2_black_100px.png)
````

#### deployment commands

yarn build --mode production (from ui dir)
scp -r dist casd:~/ddstats/ui (from ui dir)
