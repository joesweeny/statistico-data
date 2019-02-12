-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER table sportmonks_goal_event
ALTER id TYPE BIGINT;

ALTER table sportmonks_substitution_event
ALTER id TYPE BIGINT;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER table sportmonks_goal_event
ALTER id TYPE INTEGER;

ALTER table sportmonks_substitution_event
ALTER id TYPE INTEGER;