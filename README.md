# ddstats-server

![DDSTATS Server v2](/server_diagram.png)

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

#### discord-bot commands

_MAKE SURE TO CHECK FOR SURVIVAL HASH WHEN MAKING THE NOTIFICATION FUNCTION!_

- [ ] make notification function
- [x] global https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/global.js
- [x] help https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/help.js
- [x] id https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/id.js
- [x] live https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/live.js
- [x] rank https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/rank.js
- [ ] register https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/register.js
- [ ] me https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/me.js
- [x] search https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/search.js
- [x] top https://github.com/alexwilkerson/ddstats-discord-bot/blob/master/commands/top.js

#### socket-io stuff

Currently working on possibly converting a lot of the socket-io stuff to websockets. In order to stay backward-compatible with older clients, socket.io may still have to be used for client and login. Ideally, we could have the client functioning with socket.io and everything else will be running with websockets. The only thing that needs to be updated at that point would be the Discord bot which is currently receiving information via socket.io. The website could then strictly use websockets. Then at a later date, the client could be updated to use websockets as that's going to be a lot of updates to the client, I'd rather not hinder this project with that work.

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

- [x] GET player/live (gets live players from websocket)
- [x] POST api/v2/client_connect (client sends client version and receives info)
- [x] GET api/v2/motd
- [x] create table in database for motd
- [ ] GET api/v2/client/releases?page_size={int}&page_num={int}
- [ ] GET api/v2/client/download?version={string}
- [ ] create database in table for client version/filename
- [ ] GET api/v2/news?page_size={int}&page_num={int}
- [ ] create table in database for news

#### game

- [x] POST api/v2/game/submit
- [x] GET api/v2/game/top
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

### importing csvs from app.db

#### discord_users

```sql
.headers on
.mode csv
.output discord_users.csv
select * from users;
```
