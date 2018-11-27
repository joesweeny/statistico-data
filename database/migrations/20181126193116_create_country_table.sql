-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE country (
  id int NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  continent VARCHAR NOT NULL,
  iso VARCHAR NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE country