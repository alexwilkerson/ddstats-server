DROP TABLE collector;
DROP TABLE collector_run;
DROP TABLE news;
DROP TABLE release;
DROP TABLE message_of_the_day;
DROP TABLE death_type;
DROP TABLE spawnset;
DROP TABLE live;
DROP TABLE player;
DROP TABLE state;
DROP TABLE game;

CREATE TABLE IF NOT EXISTS game (
  id BIGSERIAL PRIMARY KEY NOT NULL,
  recorded TIMESTAMP WITH TIME ZONE,
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
  survival_hash TEXT,
  version TEXT,
  level_two_time DOUBLE PRECISION DEFAULT 0.0,
  level_three_time DOUBLE PRECISION DEFAULT 0.0,
  level_four_time DOUBLE PRECISION DEFAULT 0.0,
  homing_daggers_max_time DOUBLE PRECISION DEFAULT 0.0,
  enemies_alive_max_time DOUBLE PRECISION DEFAULT 0.0,
  homing_daggers_max BIGINT NOT NULL,
  enemies_alive_max BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS state (
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

CREATE INDEX game_id_idx ON state(game_id);

CREATE TABLE IF NOT EXISTS player (
  id BIGSERIAL PRIMARY KEY NOT NULL,
  last_active TIMESTAMP WITH TIME ZONE,
  player_name TEXT NOT NULL,
  rank INTEGER NOT NULL,
  game_time DOUBLE PRECISION NOT NULL,
  death_type TEXT NOT NULL,
  gems BIGINT NOT NULL,
  daggers_hit BIGINT NOT NULL,
  daggers_fired BIGINT NOT NULL,
  enemies_killed BIGINT NOT NULL,
  accuracy DOUBLE PRECISION NOT NULL,
  overall_game_time DOUBLE PRECISION NOT NULL,
  overall_average_game_time DOUBLE PRECISION NOT NULL,
  overall_deaths BIGINT NOT NULL,
  overall_gems BIGINT NOT NULL,
  overall_enemies_killed BIGINT NOT NULL,
  overall_daggers_hit BIGINT NOT NULL,
  overall_daggers_fired BIGINT NOT NULL,
  overall_accuracy DOUBLE PRECISION NOT NULL
);

CREATE TABLE IF NOT EXISTS live (
  player_id INTEGER PRIMARY KEY NOT NULL REFERENCES player(id) ON DELETE CASCADE ON UPDATE CASCADE,
  sid TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS spawnset (
  survival_hash TEXT PRIMARY KEY NOT NULL,
  spawnset_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS death_type (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS message_of_the_day (
  id SERIAL PRIMARY KEY NOT NULL,
  time_stamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  message VARCHAR(67)
);

CREATE TABLE IF NOT EXISTS discord_user (
  discord_id TEXT PRIMARY KEY,
  dd_id INTEGER NOT NULL DEFAULT 0,
  verified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS release (
  version TEXT PRIMARY KEY NOT NULL,
  time_stamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  body TEXT NOT NULL DEFAULT '',
  file_name TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS news (
  id SERIAL PRIMARY KEY,
  time_stamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  title TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS collector_run (
  id BIGSERIAL PRIMARY KEY,
  time_stamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  player_count INTEGER NOT NULL,
  new_players INTEGER NOT NULL,
  active_players INTEGER NOT NULL,
  inactive_players INTEGER NOT NULL,
  players_with_new_scores INTEGER NOT NULL,
  players_with_new_ranks INTEGER NOT NULL,
  avg_improvement_time DOUBLE PRECISION NOT NULL,
  avg_game_time_per_active_player DOUBLE PRECISION NOT NULL,
  avg_deaths_per_active_player DOUBLE PRECISION NOT NULL,
  avg_gems_per_active_player BIGINT NOT NULL,
  avg_enemies_killed_per_active_player BIGINT NOT NULL,
  avg_daggers_hit_per_active_player BIGINT NOT NULL,
  avg_daggers_fired_per_active_player BIGINT NOT NULL,
  avg_accuracy_per_active_player DOUBLE PRECISION NOT NULL,
  total_game_time DOUBLE PRECISION NOT NULL,
  deaths BIGINT NOT NULL,
  gems BIGINT NOT NULL,
  enemies_killed BIGINT NOT NULL,
  daggers_hit BIGINT NOT NULL,
  daggers_fired BIGINT NOT NULL,
  accuracy DOUBLE PRECISION NOT NULL,
);

CREATE TABLE IF NOT EXISTS collector_player (
  id INTEGER PRIMARY KEY,
  player_name TEXT NOT NULL,
  rank INTEGER NOT NULL,
  game_time DOUBLE PRECISION NOT NULL,
  death_type TEXT NOT NULL,
  gems BIGINT NOT NULL,
  daggers_hit BIGINT NOT NULL,
  daggers_fired BIGINT NOT NULL,
  enemies_killed BIGINT NOT NULL,
  overall_game_time DOUBLE PRECISION NOT NULL,
  overall_deaths BIGINT NOT NULL,
  overall_gems BIGINT NOT NULL,
  overall_enemies_killed BIGINT NOT NULL,
  overall_daggers_hit BIGINT NOT NULL,
  overall_daggers_fired BIGINT NOT NULL,
  collector_run_id BIGINT REFERENCES collector_run(id)
);

-- below are POSTGRES helper functions to make dealing with the database easier --

-- this function is used internally
CREATE OR REPLACE FUNCTION ROUND(float,int) RETURNS FLOAT AS $$
  SELECT ROUND($1::numeric,$2)::FLOAT;
$$ language SQL IMMUTABLE;

-- takes two floats and returns a float... if zero, returns 0.0
CREATE OR REPLACE FUNCTION DIVZERO(float,float) RETURNS FLOAT AS $$
  SELECT
    CASE
      WHEN $2=0 THEN 0::FLOAT
      ELSE
        ($1::FLOAT/$2::FLOAT)::FLOAT END;
$$ language SQL IMMUTABLE;