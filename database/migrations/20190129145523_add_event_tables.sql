-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_goal_event (
  id INTEGER NOT NULL PRIMARY KEY,
  team_id INTEGER NOT NULL,
  player_id INTEGER NOT NULL,
  player_assist_id INTEGER,
  minute INTEGER NOT NULL,
  score VARCHAR NOT NULL,
  created_at INTEGER NOT NULL
);

CREATE INDEX ON sportmonks_goal_event (team_id);
CREATE INDEX ON sportmonks_goal_event (player_id);
CREATE INDEX ON sportmonks_goal_event (player_assist_id);

CREATE TABLE sportmonks_substitute_event (
  id INTEGER NOT NULL PRIMARY KEY,
  team_id INTEGER NOT NULL,
  player_in_id INTEGER NOT NULL,
  player_out_id INTEGER NOT NULL,
  minute INTEGER NOT NULL,
  injured BOOLEAN,
  created_at INTEGER NOT NULL
);

CREATE INDEX ON sportmonks_substitute_event (team_id);
CREATE INDEX ON sportmonks_substitute_event (player_in_id);
CREATE INDEX ON sportmonks_substitute_event (player_out_id);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_goal_event
DROP TABLE sportmonks_substitute_event