-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_team (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  short_code VARCHAR,
  country_id INTEGER,
  venue_id INTEGER NOT NULL,
  national_team BOOLEAN NOT NULL,
  founded INTEGER,
  logo VARCHAR,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_team
