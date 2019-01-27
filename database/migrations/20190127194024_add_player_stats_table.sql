-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_player_stats (
  fixture_id INTEGER NOT NULL,
  player_id INTEGER NOT NULL,
  team_id INTEGER NOT NULL,
  position VARCHAR NOT NULL,
  formation_position INTEGER,
  substitute BOOLEAN NOT NULL,
  shots_total INTEGER,
  shots_on_goal INTEGER,
  goals_scored INTEGER,
  goals_conceded INTEGER,
  fouls_drawn INTEGER,
  fouls_committed INTEGER,
  yellow_cards INTEGER,
  red_card INTEGER,
  crosses_total INTEGER,
  crosses_accuracy INTEGER,
  passes_total INTEGER,
  passes_accuracy INTEGER,
  assists INTEGER,
  offsides INTEGER,
  saves INTEGER,
  pen_scored INTEGER,
  pen_missed INTEGER,
  pen_saved INTEGER,
  pen_committed INTEGER,
  pen_won INTEGER,
  hit_woodwork INTEGER,
  tackles INTEGER,
  blocks INTEGER,
  interceptions INTEGER,
  clearances INTEGER,
  minutes_played INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_player_stats