-- +goose Up
-- +goose StatementBegin
ALTER TABLE sportmonks_team_stats
ADD COLUMN goals INTEGER;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER table sportmonks_team_stats
DROP COLUMN goals;

-- +goose StatementEnd
