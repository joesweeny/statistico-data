-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_squad (
  season_id INT NOT NULL,
  team_id INT NOT NULL,
  player_ids INTEGER[],
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_squad