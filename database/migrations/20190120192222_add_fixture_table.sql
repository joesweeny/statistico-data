-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_fixture (
  id INTEGER NOT NULL PRIMARY KEY,
  season_id INTEGER NOT NULL,
  round_id INTEGER,
  venue_id INTEGER,
  home_team_id INTEGER NOT NULL,
  away_team_id INTEGER NOT NULL,
  referee_id INTEGER,
  date INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_fixture