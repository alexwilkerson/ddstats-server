DROP TABLE game;
DROP TABLE state;
DROP TABLE player;
DROP TABLE live;
DROP TABLE spawnset;

CREATE TABLE game (
  id BIGSERIAL PRIMARY KEY NOT NULL,
  player_id BIGINT NOT NULL,
  granularity INTEGER NOT NULL,
  game_time DOUBLE PRECISION NOT NULL,
  death_type INTEGER NOT NULL,
  gems BIGINT NOT NULL,
  homing_daggers BIGINT NOT NULL,
  daggers_fired BIGINT NOT NULL,
  daggers_hit BIGINT NOT NULL,
  enemies_alive BIGINT NOT NULL,
  enemies_killed BIGINT NOT NULL,
  time_stamp TIMESTAMP WITH TIME ZONE NOT NULL,
  replay_player_id INTEGER NOT NULL,
  survival_hash TEXT NOT NULL,
  version TEXT,
  level_two_time DOUBLE PRECISION DEFAULT 0.0,
  level_three_time DOUBLE PRECISION DEFAULT 0.0,
  level_four_time DOUBLE PRECISION DEFAULT 0.0,
  homing_daggers_max_time DOUBLE PRECISION DEFAULT 0.0,
  enemies_alive_max_time DOUBLE PRECISION DEFAULT 0.0,
  homing_daggers_max BIGINT NOT NULL,
  enemies_alive_max BIGINT NOT NULL
);

CREATE TABLE state (
  id BIGSERIAL PRIMARY KEY NOT NULL,
  game_id BIGINT NOT NULL REFERENCES game(id) ON DELETE CASCADE ON UPDATE CASCADE,
  game_time DOUBLE PRECISION NOT NULL,
  gems BIGINT NOT NULL,
  homing_daggers BIGINT NOT NULL,
  daggers_hit BIGINT NOT NULL,
  daggers_fired BIGINT NOT NULL,
  enemies_alive BIGINT NOT NULL,
  enemies_killed BIGINT NOT NULL
);

CREATE TABLE player (
  id BIGSERIAL PRIMARY KEY NOT NULL,
  username TEXT NOT NULL,
  rank INtEGER NOT NULL,
  game_time DOUBLE PRECISION NOT NULL,
  death_type TEXT NOT NULL,
  gems BIGINT NOT NULL,
  daggers_fired BIGINT NOT NULL,
  daggers_hit BIGINT NOT NULL,
  enemies_killed BIGINT NOT NULL,
  accuracy DOUBLE PRECISION NOT NULL,
  time_total DOUBLE PRECISION NOT NULL,
  deaths_total BIGINT NOT NULL,
  gems_total BIGINT NOT NULL,
  enemies_killed_total BIGINT NOT NULL,
  daggers_fired_total BIGINT NOT NULL,
  daggers_hit_total BIGINT NOT NULL,
  accuracy_total DOUBLE PRECISION NOT NULL
);

CREATE TABLE live (
  player_id INTEGER PRIMARY KEY NOT NULL REFERENCES player(id) ON DELETE CASCADE ON UPDATE CASCADE,
  sid TEXT NOT NULL
);

CREATE TABLE spawnset (
  survival_hash TEXT PRIMARY KEY NOT NULL,
  spawnset_name TEXT NOT NULL
);