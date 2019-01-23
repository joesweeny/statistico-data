-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE sportmonks_result (
    fixture_id INTEGER NOT NULL PRIMARY KEY,
    pitch_condition VARCHAR,
    home_formation VARCHAR,
    away_formation VARCHAR,
    home_score INTEGER,
    away_score INTEGER,
    home_pen_scored INTEGER,
    away_pen_scored INTEGER,
    half_time_score INTEGER,
    full_time_score INTEGER,
    extra_time_score INTEGER,
    home_league_position INTEGER,
    away_league_position INTEGER,
    minutes INTEGER,
    seconds INTEGER,
    added_time INTEGER,
    extra_time INTEGER,
    injury_time INTEGER,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sportmonks_result