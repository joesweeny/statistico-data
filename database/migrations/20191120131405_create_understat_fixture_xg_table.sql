-- +goose Up
-- +goose StatementBegin
CREATE TABLE understat_fixture_team_xg (
    id INTEGER NOT NULL PRIMARY KEY,
    sportmonks_fixture_id INTEGER NOT NULL,
    home float,
    away float,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE understat_fixture_team_xg
-- +goose StatementEnd
