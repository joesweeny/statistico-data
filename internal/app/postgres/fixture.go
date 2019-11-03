package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type FixtureRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (r *FixtureRepository) Insert(f *app.Fixture) error {
	query := `
	INSERT INTO sportmonks_fixture (id, season_id, round_id, venue_id, home_team_id, away_team_id, referee_id,
	date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.connection.Exec(
		query,
		f.ID,
		f.SeasonID,
		f.RoundID,
		f.VenueID,
		f.HomeTeamID,
		f.AwayTeamID,
		f.RefereeID,
		f.Date.Unix(),
		r.clock.Now().Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *FixtureRepository) Update(f *app.Fixture) error {
	_, err := r.ByID(uint64(f.ID))

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_fixture set season_id = $2, round_id = $3, venue_id = $4, home_team_id = $5, away_team_id = $6,
	referee_id = $7, date = $8, updated_at = $9 where id = $1`

	_, err = r.connection.Exec(
		query,
		f.ID,
		f.SeasonID,
		f.RoundID,
		f.VenueID,
		f.HomeTeamID,
		f.AwayTeamID,
		f.RefereeID,
		f.Date.Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *FixtureRepository) ByID(id uint64) (*app.Fixture, error) {
	query := `SELECT * FROM sportmonks_fixture where id = $1`
	row := r.connection.QueryRow(query, id)

	return rowToFixture(row, id)
}

func (r *FixtureRepository) IDs() ([]uint64, error) {
	t := time.Now()
	then := t.AddDate(0, 0, -1)
	query := `SELECT id FROM sportmonks_fixture where date < $1 ORDER BY id ASC`

	rows, err := r.connection.Query(query, then.Unix())

	if err != nil {
		return []uint64{}, err
	}

	return rowsToIntSlice(rows)
}

func (r *FixtureRepository) IDsBetween(from, to time.Time) ([]uint64, error) {
	query := `SELECT id FROM sportmonks_fixture where date BETWEEN $1 AND $2 ORDER BY id ASC`

	rows, err := r.connection.Query(query, from.Unix(), to.Unix())

	if err != nil {
		return []uint64{}, err
	}

	return rowsToIntSlice(rows)
}

func (r *FixtureRepository) Between(from, to time.Time) ([]app.Fixture, error) {
	query := `SELECT * FROM sportmonks_fixture where date BETWEEN $1 AND $2 ORDER BY id ASC`

	rows, err := r.connection.Query(query, from.Unix(), to.Unix())

	if err != nil {
		return []app.Fixture{}, err
	}

	return rowsToFixtureSlice(rows)
}

func (r *FixtureRepository) ByTeamID(id uint64, limit int32, before time.Time) ([]app.Fixture, error) {
	query := `SELECT * FROM sportmonks_fixture WHERE date < $2 AND (home_team_id = $1 OR away_team_id = $1)
	ORDER BY date DESC LIMIT $3`

	rows, err := r.connection.Query(query, id, before.Unix(), limit)

	if err != nil {
		return []app.Fixture{}, err
	}

	return rowsToFixtureSlice(rows)
}

func (r *FixtureRepository) BySeasonID(id uint64, before time.Time) ([]app.Fixture, error) {
	query := `SELECT * FROM sportmonks_fixture WHERE season_id = $1 and date < $2 ORDER BY date ASC, id ASC`

	rows, err := r.connection.Query(query, id, before.Unix())

	if err != nil {
		return []app.Fixture{}, err
	}

	return rowsToFixtureSlice(rows)
}

func (r *FixtureRepository) ByHomeAndAwayTeam(homeTeamId, awayTeamId uint64, limit uint32, before time.Time) ([]app.Fixture, error) {
	query := `SELECT * FROM sportmonks_fixture WHERE home_team_id = $1 and away_team_id = $2 and date < $3
	ORDER BY date DESC LIMIT $4`

	rows, err := r.connection.Query(query, homeTeamId, awayTeamId, before.Unix(), limit)

	if err != nil {
		return []app.Fixture{}, err
	}

	return rowsToFixtureSlice(rows)
}

func rowsToIntSlice(rows *sql.Rows) ([]uint64, error) {
	defer rows.Close()

	var id uint64
	var ids []uint64

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func rowsToFixtureSlice(rows *sql.Rows) ([]app.Fixture, error) {
	defer rows.Close()

	var date int64
	var created int64
	var updated int64
	var fixtures []app.Fixture
	var f app.Fixture

	for rows.Next() {
		err := rows.Scan(
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
			return fixtures, err
		}

		f.Date = time.Unix(date, 0)
		f.CreatedAt = time.Unix(created, 0)
		f.UpdatedAt = time.Unix(updated, 0)

		fixtures = append(fixtures, f)
	}

	return fixtures, nil
}

func rowToFixture(r *sql.Row, id uint64) (*app.Fixture, error) {
	var date int64
	var created int64
	var updated int64

	f := app.Fixture{}

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
		return &f, fmt.Errorf("fixture with ID %d does not exist", id)
	}

	f.Date = time.Unix(date, 0)
	f.CreatedAt = time.Unix(created, 0)
	f.UpdatedAt = time.Unix(updated, 0)

	return &f, nil
}

func NewFixtureRepository(connection *sql.DB, clock clockwork.Clock) *FixtureRepository {
	return &FixtureRepository{connection: connection, clock: clock}
}
