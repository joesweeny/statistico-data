-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE sportmonks_player_stats
ALTER COLUMN position DROP NOT NULL;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE sportmonks_player_stats
ALTER COLUMN position SET NOT NULL;