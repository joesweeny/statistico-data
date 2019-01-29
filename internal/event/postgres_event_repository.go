package event

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
)

type PostgresEventRepository struct {
	Connection *sql.DB
}

func (p *PostgresEventRepository) InsertGoalEvent(m *model.GoalEvent) error {
	query := `
	INSERT INTO sportmonks_goal_event (id, team_id, player_id, player_assist_id, minute, score, created_at) VALUES
	($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.Connection.Exec(
		query,
		m.ID,
		m.TeamID,
		m.PlayerID,
		m.PlayerAssistID,
		m.Minute,
		m.Score,
		m.CreatedAt.Unix(),
	)

	return err
}

func (p *PostgresEventRepository) InsertSubstitutionEvent(m *model.SubstitutionEvent) error {
	query := `
	INSERT INTO sportmonks_goal_event (id, team_id, player_in_id, player_out_id, minute, injured, created_at) VALUES
	($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.Connection.Exec(
		query,
		m.ID,
		m.TeamID,
		m.PlayerInID,
		m.PlayerOutID,
		m.Minute,
		m.Injured,
		m.CreatedAt.Unix(),
	)

	return err
}