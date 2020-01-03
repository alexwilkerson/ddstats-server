# ddstats-api

### TODO: endpoints

#### index

- [ ] GET user/live
- [ ] GET client/version/latest
- [ ] GET news/recent?num={int}

#### game

- [ ] GET api/v2/game/top
- [x] GET api/v2/game/recent?pagesize={int}&pagenum={int}
- [x] GET api/v2/game?id={int} (info)
- [x] GET api/v2/game/all?id={int}
- [x] GET api/v2/game/gems?id={int}
- [x] GET api/v2/game/homing-daggers?id={int}
- [x] GET api/v2/game/accuracy?id={int}
- [x] GET api/v2/game/enemies-alive?id={int}
- [x] GET api/v2/game/enemies-killed?id={int}

#### player

- [x] GET api/v2/player/all?pagesize={int}&pagenum={int} (list all users)
- [x] GET api/v2/player?id={int}

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
select id, username, rank, game_time, death_type, gems, daggers_fired, daggers_hit, enemies_killed, accuracy, time_total as overall_time, deaths_total as overall_deaths, gems_total as overall_gems, enemies_killed_total as overall_enemies_killed, daggers_fired_total as overall_daggers_fired, daggers_hit_total as overall_daggers_hit, accuracy_total as overall_accuracy from user;
```

#### live

```sql
.headers on
.mode csv
.output live.csv
select * from live;
```
