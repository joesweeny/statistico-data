-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_round (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  season_id INTEGER NOT NULL,
  start_date INTEGER NOT NULL,
  end_date INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_round