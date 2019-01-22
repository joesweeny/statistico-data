-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_venue (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  surface VARCHAR,
  address VARCHAR,
  city VARCHAR,
  capacity INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_venue