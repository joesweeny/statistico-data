package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/errors"
	"time"
)

type TeamRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (r *TeamRepository) Insert(t *app.Team) error {
	query := `
	INSERT INTO sportmonks_team (id, name, short_code, country_id, venue_id, national_team, founded, logo, created_at,
	updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.connection.Exec(
		query,
		t.ID,
		t.Name,
		t.ShortCode,
		t.CountryID,
		t.VenueID,
		t.NationalTeam,
		t.Founded,
		t.Logo,
		r.clock.Now().Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r TeamRepository) Update(m *app.Team) error {
	_, err := r.ByID(m.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_team set name = $2, short_code = $3, country_id = $4, venue_id = $5, national_team = $6,
	founded = $7, logo = $8, updated_at = $9 where id = $1`

	_, err = r.connection.Exec(
		query,
		m.ID,
		m.Name,
		m.ShortCode,
		m.CountryID,
		m.VenueID,
		m.NationalTeam,
		m.Founded,
		m.Logo,
		r.clock.Now().Unix(),
	)

	return err
}

func (r TeamRepository) ByID(id uint64) (*app.Team, error) {
	query := `SELECT * FROM sportmonks_team where id = $1`
	row := r.connection.QueryRow(query, id)

	return rowToTeam(row)
}

func (r TeamRepository) BySeasonId(id uint64) ([]*app.Team, error) {
	builder := r.queryBuilder()

	sub := builder.Select("*").
		Prefix("IN").
		From("sportmonks_fixture").
		Distinct().
		Options("home_team_id").
		Where(sq.Eq{"season_id": id})

	query := builder.Select("sportmonks_team.*").
		From("sportmonks_team").
		Where(sub).
		OrderBy("sportmonks_team.name ASC")

	rows, err := query.Query()

	if err != nil {
		return []*app.Team{}, err
	}

	var created int64
	var updated int64
	var teams []*app.Team
	var team app.Team

	for rows.Next() {
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.ShortCode,
			&team.CountryID,
			&team.VenueID,
			&team.NationalTeam,
			&team.Founded,
			&team.Logo,
			&created,
			&updated,
		)

		if err != nil {
			return []*app.Team{}, err
		}

		team.CreatedAt = time.Unix(created, 0)
		team.UpdatedAt = time.Unix(updated, 0)

		teams = append(teams, &team)
	}

	return teams, nil
}

func rowToTeam(r *sql.Row) (*app.Team, error) {
	var created int64
	var updated int64

	var t = app.Team{}

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
		if err == sql.ErrNoRows {
			return nil, errors.ErrorNotFound
		}

		return nil, err
	}

	t.CreatedAt = time.Unix(created, 0)
	t.UpdatedAt = time.Unix(updated, 0)

	return &t, nil
}

func (r TeamRepository) queryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(r.connection)
}

func NewTeamRepository(connection *sql.DB, clock clockwork.Clock) *TeamRepository {
	return &TeamRepository{connection: connection, clock: clock}
}
