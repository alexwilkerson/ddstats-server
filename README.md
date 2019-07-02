### importing csvs

#### state.csv
```sql
.headers on
.mode csv
.output state.csv
delete from state where game_id not in (select id from game where id is not null);
select id, game_id, printf("%.6f", game_time) as game_time, gems, homing_daggers, daggers_hit, daggers_fired, enemies_alive, enemies_killed from state;
```