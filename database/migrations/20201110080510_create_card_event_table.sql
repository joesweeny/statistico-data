-- +goose Up
-- +goose StatementBegin
CREATE TABLE sportmonks_card_event (
  id BIGINT NOT NULL PRIMARY KEY,
  team_id INTEGER NOT NULL,
  fixture_id INTEGER NOT NULL,
  type VARCHAR NOT NULL,
  player_id INTEGER NOT NULL,
  minute INTEGER NOT NULL,
  reason VARCHAR,
  created_at INTEGER NOT NULL
);

CREATE INDEX ON sportmonks_card_event (fixture_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sportmonks_card_event
-- +goose StatementEnd
