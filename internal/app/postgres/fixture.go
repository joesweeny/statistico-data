package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
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

	query := `UPDATE sportmonks_fixture set season_id = $2, round_id = $3, venue_id = $4, home_team_id = $5, away_team_id = $6,
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

func (r *FixtureRepository) ByTeamID(id uint64, query app.FixtureFilterQuery) ([]app.Fixture, error) {
	builder := r.queryBuilder()

	q := builder.Select("sportmonks_fixture.*").From("sportmonks_fixture")

	if query.Limit != nil {
		q = q.Limit(*query.Limit)
	}

	if query.DateBefore != nil {
		q = q.Where(sq.Lt{"date": query.DateBefore.Unix()})
	}

	if query.Venue == nil {
		q = q.Where(sq.Or{
			sq.Eq{"home_team_id": id},
			sq.Eq{"away_team_id": id},
		})
	} else {
		if *query.Venue == "home" {
			q = q.Where(sq.Eq{"home_team_id": id})
		}

		if *query.Venue == "away" {
			q = q.Where(sq.Eq{"away_team_id": id})
		}
	}

	if query.SortBy != nil && *query.SortBy == "date_asc" {
		q = q.OrderBy("date ASC")
	}

	if query.SortBy != nil && *query.SortBy == "date_desc" {
		q = q.OrderBy("date DESC")
	}

	if query.SortBy == nil {
		q = q.OrderBy("date DESC")
	}

	rows, err := q.Query()

	if err != nil {
		return []app.Fixture{}, err
	}

	return rowsToFixtureSlice(rows)
}

func (r *FixtureRepository) Get(q app.FixtureRepositoryQuery) ([]app.Fixture, error) {
	builder := r.queryBuilder()

	query := builder.Select("sportmonks_fixture.*").From("sportmonks_fixture")

	rows, err := buildQuery(query, q).Query()

	if err != nil {
		return []app.Fixture{}, err
	}

	return rowsToFixtureSlice(rows)
}

func (r *FixtureRepository) GetIDs(q app.FixtureRepositoryQuery) ([]uint64, error) {
	builder := r.queryBuilder()

	query := builder.Select("id").From("sportmonks_fixture")

	rows, err := buildQuery(query, q).Query()

	if err != nil {
		return []uint64{}, err
	}

	return rowsToIntSlice(rows)
}

func buildQuery(b sq.SelectBuilder, q app.FixtureRepositoryQuery) sq.SelectBuilder {
	if q.SeasonID != nil {
		b = b.Where(sq.Eq{"season_id": q.SeasonID})
	}

	if q.HomeTeamID != nil {
		b = b.Where(sq.Eq{"home_team_id": q.HomeTeamID})
	}

	if q.AwayTeamID != nil {
		b = b.Where(sq.Eq{"away_team_id": q.AwayTeamID})
	}

	if q.HomeTeamNameLike != nil {
		nested := sq.Select("*").From("sportmonks_team").Where(sq.Like{"sportmonks_team.name": *q.HomeTeamNameLike + "%"})

		b = b.JoinClause(nested.Prefix("JOIN (").Suffix(") t1 ON sportmonks_fixture.home_team_id = t1.id"))
	}

	if q.AwayTeamNameLike != nil {
		nested := sq.Select("*").From("sportmonks_team").Where(sq.Like{"sportmonks_team.name": *q.AwayTeamNameLike + "%"})

		b = b.JoinClause(nested.Prefix("JOIN (").Suffix(") t2 ON sportmonks_fixture.away_team_id = t2.id"))
	}

	if q.DateFrom != nil {
		b = b.Where(sq.GtOrEq{"date": q.DateFrom.Unix()})
	}

	if q.DateTo != nil {
		b = b.Where(sq.LtOrEq{"date": q.DateTo.Unix()})
	}

	if q.Limit != nil {
		b = b.Limit(*q.Limit)
	}

	if q.SortBy != nil && *q.SortBy == "date_asc" {
		b = b.OrderBy("date ASC")
	}

	if q.SortBy != nil && *q.SortBy == "date_desc" {
		b = b.OrderBy("date DESC")
	}

	if q.SortBy == nil {
		b = b.OrderBy("date ASC")
	}

	return b
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

func (r *FixtureRepository) queryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(r.connection)
}
