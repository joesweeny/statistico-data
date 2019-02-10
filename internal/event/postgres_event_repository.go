package event

import (
	"database/sql"
	"errors"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresEventRepository struct {
	Connection *sql.DB
}

func (p *PostgresEventRepository) InsertGoalEvent(m *model.GoalEvent) error {
	query := `
	INSERT INTO sportmonks_goal_event (id, team_id, player_id, player_assist_id, minute, score, created_at, fixture_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := p.Connection.Exec(
		query,
		m.ID,
		m.TeamID,
		m.PlayerID,
		m.PlayerAssistID,
		m.Minute,
		m.Score,
		m.CreatedAt.Unix(),
		m.FixtureID,
	)

	return err
}

func (p *PostgresEventRepository) InsertSubstitutionEvent(m *model.SubstitutionEvent) error {
	query := `
	INSERT INTO sportmonks_substitution_event (id, team_id, player_in_id, player_out_id, minute, injured, created_at, 
	fixture_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := p.Connection.Exec(
		query,
		m.ID,
		m.TeamID,
		m.PlayerInID,
		m.PlayerOutID,
		m.Minute,
		m.Injured,
		m.CreatedAt.Unix(),
		m.FixtureID,
	)

	return err
}

func (p *PostgresEventRepository) GoalEventById(id int) (*model.GoalEvent, error) {
	query := `SELECT * FROM sportmonks_goal_event WHERE id = $1`

	m := model.GoalEvent{}

	var created int64

	row := p.Connection.QueryRow(query, id)

	err := row.Scan(&m.ID, &m.TeamID, &m.PlayerID, &m.PlayerAssistID, &m.Minute, &m.Score, &created, &m.FixtureID)

	if err != nil {
		return &m, ErrNotFound
	}

	m.CreatedAt = time.Unix(created, 0)

	return &m, nil
}

func (p *PostgresEventRepository) SubstitutionEventById(id int) (*model.SubstitutionEvent, error) {
	query := `SELECT * FROM sportmonks_substitution_event WHERE id = $1`

	m := model.SubstitutionEvent{}

	var created int64

	row := p.Connection.QueryRow(query, id)

	err := row.Scan(&m.ID, &m.TeamID, &m.PlayerInID, &m.PlayerOutID, &m.Minute, &m.Injured, &created, &m.FixtureID)

	if err != nil {
		return &m, ErrNotFound
	}

	m.CreatedAt = time.Unix(created, 0)

	return &m, nil
}
