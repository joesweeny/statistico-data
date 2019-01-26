-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_player (
  id INTEGER NOT NULL PRIMARY KEY,
  team_id INTEGER,
  country_id INTEGER NOT NULL,
  first_name VARCHAR NOT NULL,
  last_name VARCHAR NOT NULL,
  birth_place VARCHAR,
  date_of_birth DATE,
  position_id INTEGER NOT NULL,
  image VARCHAR,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_player