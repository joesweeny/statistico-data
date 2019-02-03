package team

import (
	"database/sql"
	"errors"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresTeamRepository struct {
	Connection *sql.DB
}

func (p *PostgresTeamRepository) Insert(t *model.Team) error {
	query := `
	INSERT INTO sportmonks_team (id, name, short_code, country_id, venue_id, national_team, founded, logo, created_at,
	updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.Connection.Exec(
		query,
		t.ID,
		t.Name,
		t.ShortCode,
		t.CountryID,
		t.VenueID,
		t.NationalTeam,
		t.Founded,
		t.Logo,
		t.CreatedAt.Unix(),
		t.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresTeamRepository) Update(m *model.Team) error {
	_, err := p.GetById(m.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_team set name = $2, short_code = $3, country_id = $4, venue_id = $5, national_team = $6,
	founded = $7, logo = $8, updated_at = $9 where id = $1`

	_, err = p.Connection.Exec(
		query,
		m.ID,
		m.Name,
		m.ShortCode,
		m.CountryID,
		m.VenueID,
		m.NationalTeam,
		m.Founded,
		m.Logo,
		m.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresTeamRepository) GetById(id int) (*model.Team, error) {
	query := `SELECT * FROM sportmonks_team where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToTeam(row)
}

func rowToTeam(r *sql.Row) (*model.Team, error) {
	var created int64
	var updated int64

	t := model.Team{}

	err := r.Scan(
		&t.ID,
		&t.Name,
		&t.ShortCode,
		&t.CountryID,
		&t.VenueID,
		&t.NationalTeam,
		&t.Founded,
		&t.Logo,
		&created,
		&updated,
	)

	if err != nil {
		return &t, ErrNotFound
	}

	t.CreatedAt = time.Unix(created, 0)
	t.UpdatedAt = time.Unix(updated, 0)

	return &t, nil
}
