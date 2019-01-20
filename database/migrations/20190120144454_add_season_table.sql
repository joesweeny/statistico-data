-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_season (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  league_id INTEGER NOT NULL,
  is_current BOOLEAN NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_season