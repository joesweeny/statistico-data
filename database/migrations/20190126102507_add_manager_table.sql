-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_manager (
  id INTEGER NOT NULL PRIMARY KEY,
  team_id INTEGER,
  first_name VARCHAR NOT NULL,
  last_name VARCHAR NOT NULL,
  birth_country VARCHAR NOT NULL,
  birthplace VARCHAR NOT NULL,
  nationality VARCHAR NOT NULL,
  date_of_birth DATE NOT NULL,
  image VARCHAR,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_manager