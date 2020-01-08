# ddstats-api

## Automatically restarting server during dev

Go 1.11+ installation required\
Install reflex: `go get github.com/cespare/reflex`

### Things Left to Do...

#### socket-io stuff

- [ ] clear live table func
- [ ] create get_live_users in /ddstats-bot namespace https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L344-L354
- [ ] create on get_status in namespace /stats https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L331-L337
- [ ] create game_submitted in namespace /stats https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L311-L328
- [ ] create submit in namespace /stats https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L279-L308
- [ ] create status_update https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L217-L222
- [ ] create submit https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L224-L253
- [ ] create game_submitted https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L255-L272
- [ ] create disconnect https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L197-L206
- [ ] create login https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L175-L194
- [ ] create join in namespace /user-page https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L160-L172
- [ ] create disconnect in namespace /user-page https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L138-L152
- [ ] create a userList and playerDict https://github.com/alexwilkerson/ddstats/blob/master/socketio_main.py#L116-L118

#### submission stuff

- [x] finish the logic behind game submissions and filtering (i think this is done)

### TODO: endpoints

#### index

- [x] POST api/v2/client_connect (client sends client version and receives info)
- [x] GET api/v2/motd
- [x] create table in database for motd
- [ ] GET api/v2/client/releases?pagesize={int}&pagenum={int}
- [ ] GET api/v2/client/download?version={string}
- [ ] create database in table for client version/filename
- [ ] GET api/v2/news?pagesize={int}&pagenum={int}
- [ ] create table in database for news

- [ ] GET user/live (idk about this one?)

#### game

- [x] POST api/v2/game/submit
- [x] GET api/v2/game/top
- [x] GET api/v2/game/recent?pagesize={int}&pagenum={int} (optional: playerid={int} will give most recent for given player)
- [x] GET api/v2/game?id={int} (info)
- [x] GET api/v2/game/all?id={int}
- [x] GET api/v2/game/gems?id={int}
- [x] GET api/v2/game/homing-daggers?id={int}
- [x] GET api/v2/game/accuracy?id={int}
- [x] GET api/v2/game/enemies-alive?id={int}
- [x] GET api/v2/game/enemies-killed?id={int}

#### player

- [x] GET api/v2/player/update?id={int} (updates or creates user in database from dd backend)
- [x] GET api/v2/player/all?pagesize={int}&pagenum={int} (list all users)
- [x] GET api/v2/player?id={int}
- [x] ~GET api/v2/player/games?id={int}pagesize={int}&pagenum={int} (paginated list of games by player id)~
      (see api/v2/game/recent)

#### dd backend

- [x] GET api/v2/ddapi/get_user_by_rank?rank={int}
- [x] GET api/v2/ddapi/get_user_by_id?id={int}
- [x] GET api/v2/ddapi/user_search?user={string}
- [x] GET api/v2/ddapi/get_scores?offset={int}&limit={int}

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
select id, username as player_name, rank, game_time, death_type, gems, daggers_fired, daggers_hit, enemies_killed, accuracy, time_total as overall_time, deaths_total as overall_deaths, gems_total as overall_gems, enemies_killed_total as overall_enemies_killed, daggers_fired_total as overall_daggers_fired, daggers_hit_total as overall_daggers_hit, accuracy_total as overall_accuracy from user;
```

#### live

```sql
.headers on
.mode csv
.output live.csv
select * from live;
```
