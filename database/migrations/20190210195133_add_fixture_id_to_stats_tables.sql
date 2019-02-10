-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE sportmonks_goal_event
ADD COLUMN fixture_id INTEGER NOT NULL;

CREATE INDEX ON sportmonks_goal_event (fixture_id);

ALTER TABLE sportmonks_substitution_event
ADD COLUMN fixture_id INTEGER NOT NULL;

CREATE INDEX ON sportmonks_substitution_event (fixture_id);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE sportmonks_goal_event
DROP COLUMN fixture_id;

ALTER TABLE sportmonks_substitution_event
DROP COLUMN fixture_id;