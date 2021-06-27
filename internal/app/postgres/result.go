package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-football-data/internal/app"
	"time"
)

type ResultRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (p *ResultRepository) Insert(r *app.Result) error {
	query := `
	INSERT INTO sportmonks_result (fixture_id, pitch_condition, home_formation, away_formation, home_score, away_score,
	home_pen_score, away_pen_score, half_time_score, full_time_score, extra_time_score, home_league_position,
	away_league_position, minutes, added_time, extra_time, injury_time, created_at, updated_at) VALUES ($1, $2, 
	$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	_, err := p.connection.Exec(
		query,
		r.FixtureID,
		r.PitchCondition,
		r.HomeFormation,
		r.AwayFormation,
		r.HomeScore,
		r.AwayScore,
		r.HomePenScore,
		r.AwayPenScore,
		r.HalfTimeScore,
		r.FullTimeScore,
		r.ExtraTimeScore,
		r.HomeLeaguePosition,
		r.AwayLeaguePosition,
		r.Minutes,
		r.AddedTime,
		r.ExtraTime,
		r.InjuryTime,
		p.clock.Now().Unix(),
		p.clock.Now().Unix(),
	)

	return err
}

func (p *ResultRepository) Update(r *app.Result) error {
	_, err := p.ByFixtureID(r.FixtureID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_result SET pitch_condition = $2, home_formation = $3, away_formation = $4, home_score = $5, 
    away_score = $6, home_pen_score = $7, away_pen_score = $8, half_time_score = $9, full_time_score = $10, 
	extra_time_score = $11, home_league_position = $12, away_league_position = $13, minutes = $14, 
	added_time = $15, extra_time = $16, injury_time = $17, updated_at = $18 WHERE fixture_id = $1`

	_, err = p.connection.Exec(
		query,
		r.FixtureID,
		r.PitchCondition,
		r.HomeFormation,
		r.AwayFormation,
		r.HomeScore,
		r.AwayScore,
		r.HomePenScore,
		r.AwayPenScore,
		r.HalfTimeScore,
		r.FullTimeScore,
		r.ExtraTimeScore,
		r.HomeLeaguePosition,
		r.AwayLeaguePosition,
		r.Minutes,
		r.AddedTime,
		r.ExtraTime,
		r.InjuryTime,
		p.clock.Now().Unix(),
	)

	return err
}

func (p *ResultRepository) ByFixtureID(id uint64) (*app.Result, error) {
	query := `SELECT * FROM sportmonks_result where fixture_id = $1`
	row := p.connection.QueryRow(query, id)

	return rowToResult(row, id)
}

func rowToResult(r *sql.Row, id uint64) (*app.Result, error) {
	var created int64
	var updated int64

	var m = app.Result{}

	err := r.Scan(
		&m.FixtureID,
		&m.PitchCondition,
		&m.HomeFormation,
		&m.AwayFormation,
		&m.HomeScore,
		&m.AwayScore,
		&m.HomePenScore,
		&m.AwayPenScore,
		&m.HalfTimeScore,
		&m.FullTimeScore,
		&m.ExtraTimeScore,
		&m.HomeLeaguePosition,
		&m.AwayLeaguePosition,
		&m.Minutes,
		&m.AddedTime,
		&m.ExtraTime,
		&m.InjuryTime,
		&created,
		&updated,
	)

	if err != nil {
		return &m, fmt.Errorf("result with ID %d does not exist", id)
	}

	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}

func NewResultRepository(connection *sql.DB, clock clockwork.Clock) *ResultRepository {
	return &ResultRepository{connection: connection, clock: clock}
}
