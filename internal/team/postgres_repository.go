package team

import (
	"errors"
	"database/sql"
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

func (p *PostgresTeamRepository) GetById(id int) (*model.Team, error) {
	query := `SELECT * FROM sportmonks_team where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToTeam(row)
}

func rowToTeam(r *sql.Row) (*model.Team, error) {
	var id int
	var name string
	var short *string
	var country *int
	var venue int
	var national bool
	var founded *int
	var logo *string
	var created int64
	var updated int64

	t := model.Team{}

	if err := r.Scan(&id, &name, &short, &country, &venue, &national, &founded, &logo, &created, &updated); err != nil {
		return &t, ErrNotFound
	}

	t.ID = id
	t.Name = name
	t.ShortCode = short
	t.CountryID = country
	t.VenueID = venue
	t.NationalTeam = national
	t.Founded = founded
	t.Logo = logo
	t.CreatedAt = time.Unix(created, 0)
	t.UpdatedAt = time.Unix(updated, 0)

	return &t, nil
}
