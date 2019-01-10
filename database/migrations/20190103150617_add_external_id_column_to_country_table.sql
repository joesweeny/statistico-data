-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE country
ALTER COLUMN id TYPE VARCHAR,
ALTER COLUMN id SET NOT NULL;

ALTER TABLE country
ADD external_id int NOT NULL;

CREATE INDEX country_external_id
ON country (external_id);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE country
DROP COLUMN external_id