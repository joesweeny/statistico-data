-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE country
RENAME TO sportmonks_country;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE sportmonks_country
RENAME TO country;