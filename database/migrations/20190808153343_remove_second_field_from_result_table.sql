-- +goose Up
-- +goose StatementBegin
ALTER table sportmonks_result DROP COLUMN seconds;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER table sportmonks_result COLUMN seconds;
-- +goose StatementEnd
