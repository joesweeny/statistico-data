package fixture

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresFixtureRepository struct {
	Connection *sql.DB
}

func (p *PostgresFixtureRepository) Insert(f *model.Fixture) error {
	query := `
	INSERT INTO sportmonks_fixture (id, season_id, round_id, venue_id, home_team_id, away_team_id, referee_id,
	date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.Connection.Exec(
		query,
		f.ID,
		f.SeasonID,
		f.RoundID,
		f.VenueID,
		f.HomeTeamID,
		f.AwayTeamID,
		f.RefereeID,
		f.Date.Unix(),
		f.CreatedAt.Unix(),
		f.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresFixtureRepository) Update(f *model.Fixture) error {
	_, err := p.GetById(f.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_fixture set season_id = $2, round_id = $3, venue_id = $4, home_team_id = $5, away_team_id = $6,
	referee_id = $7, date = $8, updated_at = $9 where id = $1`

	_, err = p.Connection.Exec(
		query,
		f.ID,
		f.SeasonID,
		f.RoundID,
		f.VenueID,
		f.HomeTeamID,
		f.AwayTeamID,
		f.RefereeID,
		f.Date.Unix(),
		f.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresFixtureRepository) GetById(id int) (*model.Fixture, error) {
	query := `SELECT * FROM sportmonks_fixture where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToFixture(row)
}

func (p *PostgresFixtureRepository) Ids() ([]int, error) {
	t := time.Now()
	then := t.AddDate(0, 0, -1)
	query := `SELECT id FROM sportmonks_fixture where date < $1 ORDER BY id ASC`

	rows, err := p.Connection.Query(query, then.Unix())

	if err != nil {
		return []int{}, err
	}

	return rowsToSlice(rows)
}

func (p *PostgresFixtureRepository) IdsBetween(from, to time.Time) ([]int, error) {
	query := `SELECT id FROM sportmonks_fixture where date BETWEEN $1 AND $2 ORDER BY id ASC`

	rows, err := p.Connection.Query(query, from.Unix(), to.Unix())

	if err != nil {
		return []int{}, err
	}

	return rowsToSlice(rows)
}

func rowsToSlice(rows *sql.Rows) ([]int, error) {
	defer rows.Close()

	var id int
	var ids []int

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func rowToFixture(r *sql.Row) (*model.Fixture, error) {
	var date int64
	var created int64
	var updated int64

	f := model.Fixture{}

	err := r.Scan(
		&f.ID,
		&f.SeasonID,
		&f.RoundID,
		&f.VenueID,
		&f.HomeTeamID,
		&f.AwayTeamID,
		&f.RefereeID,
		&date,
		&created,
		&updated,
	)

	if err != nil {
		return &f, ErrNotFound
	}

	f.Date = time.Unix(date, 0)
	f.CreatedAt = time.Unix(created, 0)
	f.UpdatedAt = time.Unix(updated, 0)

	return &f, nil
}
