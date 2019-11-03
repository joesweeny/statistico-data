package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type RoundRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (p *RoundRepository) Insert(r *app.Round) error {
	query := `
	INSERT INTO sportmonks_round (id, name, season_id, start_date, end_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.connection.Exec(
		query,
		r.ID,
		r.Name,
		r.SeasonID,
		r.StartDate.Unix(),
		r.EndDate.Unix(),
		p.clock.Now().Unix(),
		p.clock.Now().Unix(),
	)

	return err
}

func (p *RoundRepository) ByID(id uint64) (*app.Round, error) {
	query := `SELECT * FROM sportmonks_round where id = $1`
	row := p.connection.QueryRow(query, id)

	return rowToRound(row)
}

func (p *RoundRepository) Update(r *app.Round) error {
	_, err := p.ByID(r.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_round set name = $2, season_id = $3, start_date = $4, end_date = $5, updated_at = $6
	where id = $1`

	_, err = p.connection.Exec(
		query,
		r.ID,
		r.Name,
		r.SeasonID,
		r.StartDate.Unix(),
		r.EndDate.Unix(),
		p.clock.Now().Unix(),
	)

	return err
}

func rowToRound(r *sql.Row) (*app.Round, error) {
	var start int64
	var end int64
	var created int64
	var updated int64

	var m = app.Round{}
	err := r.Scan(&m.ID, &m.Name, &m.SeasonID, &start, &end, &created, &updated)

	if err != nil {
		return &m, errors.New(fmt.Sprintf("Season with ID %d does not exist", m.ID))
	}

	m.StartDate = time.Unix(start, 0)
	m.EndDate = time.Unix(end, 0)
	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}

func NewRoundRepository(connection *sql.DB, clock clockwork.Clock) *RoundRepository {
	return &RoundRepository{connection: connection, clock: clock}
}