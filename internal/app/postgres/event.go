package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type EventRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (e *EventRepository) InsertGoalEvent(g *app.GoalEvent) error {
	query := `
	INSERT INTO sportmonks_goal_event (id, team_id, player_id, player_assist_id, minute, score, created_at, fixture_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := e.connection.Exec(
		query,
		g.ID,
		g.TeamID,
		g.PlayerID,
		g.PlayerAssistID,
		g.Minute,
		g.Score,
		e.clock.Now().Unix(),
		g.FixtureID,
	)

	return err
}

func (e *EventRepository) InsertSubstitutionEvent(s *app.SubstitutionEvent) error {
	query := `
	INSERT INTO sportmonks_substitution_event (id, team_id, player_in_id, player_out_id, minute, injured, created_at, 
	fixture_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := e.connection.Exec(
		query,
		s.ID,
		s.TeamID,
		s.PlayerInID,
		s.PlayerOutID,
		s.Minute,
		s.Injured,
		e.clock.Now().Unix(),
		s.FixtureID,
	)

	return err
}

func (e *EventRepository) GoalEventByID(id uint64) (*app.GoalEvent, error) {
	query := `SELECT * FROM sportmonks_goal_event WHERE id = $1`

	var g = app.GoalEvent{}
	var created int64

	row := e.connection.QueryRow(query, id)

	err := row.Scan(&g.ID, &g.TeamID, &g.PlayerID, &g.PlayerAssistID, &g.Minute, &g.Score, &created, &g.FixtureID)

	if err != nil {
		return &g, fmt.Errorf("goal event with ID %d does not exist", id)
	}

	g.CreatedAt = time.Unix(created, 0)

	return &g, nil
}

func (g *EventRepository) SubstitutionEventByID(id uint64) (*app.SubstitutionEvent, error) {
	query := `SELECT * FROM sportmonks_substitution_event WHERE id = $1`

	var s = app.SubstitutionEvent{}
	var created int64

	row := g.connection.QueryRow(query, id)

	err := row.Scan(&s.ID, &s.TeamID, &s.PlayerInID, &s.PlayerOutID, &s.Minute, &s.Injured, &created, &s.FixtureID)

	if err != nil {
		return &s, fmt.Errorf("substitution event with ID %d does not exist", id)
	}

	s.CreatedAt = time.Unix(created, 0)

	return &s, nil
}

func NewEventRepository(connection *sql.DB, clock clockwork.Clock) *EventRepository {
	return &EventRepository{connection: connection, clock: clock}
}
