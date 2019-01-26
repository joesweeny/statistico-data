-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_team_stats (
  fixture_id INTEGER NOT NULL,
  team_id INTEGER NOT NULL,
  shots_total INTEGER,
  shots_on_goal INTEGER,
  shots_off_goal INTEGER,
  shots_blocked INTEGER,
  shots_inside_box INTEGER,
  shots_outside_box INTEGER,
  passes_total INTEGER,
  passes_accuracy INTEGER,
  passes_percentage INTEGER,
  attacks_total INTEGER,
  attacks_dangerous INTEGER,
  fouls INTEGER,
  corners INTEGER,
  offsides INTEGER,
  possession INTEGER,
  yellow_cards INTEGER,
  red_cards INTEGER,
  saves INTEGER,
  substitutions INTEGER,
  goal_kicks INTEGER,
  goal_attempts INTEGER,
  free_kicks INTEGER,
  throw_ins INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_team_stats