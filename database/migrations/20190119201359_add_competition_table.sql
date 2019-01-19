-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_competition (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  country_id INTEGER NOT NULL,
  is_cup BOOLEAN NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_competition