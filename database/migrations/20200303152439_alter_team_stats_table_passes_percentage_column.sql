-- +goose Up
-- +goose StatementBegin
ALTER table sportmonks_team_stats
ALTER passes_percentage TYPE DECIMAL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER table sportmonks_team_stats
ALTER passes_percentage TYPE INTEGER ;
-- +goose StatementEnd
