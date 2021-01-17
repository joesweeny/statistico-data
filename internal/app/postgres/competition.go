package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type CompetitionRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (r *CompetitionRepository) Insert(c *app.Competition) error {
	query := `
	INSERT INTO sportmonks_competition (id, name, country_id, is_cup, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.connection.Exec(
		query,
		c.ID,
		c.Name,
		c.CountryID,
		c.IsCup,
		r.clock.Now().Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *CompetitionRepository) Update(c *app.Competition) error {
	_, err := r.ByID(c.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_competition set name = $2, is_cup = $3, updated_at = $4 where id = $1`

	_, err = r.connection.Exec(
		query,
		c.ID,
		c.Name,
		c.IsCup,
		r.clock.Now().Unix(),
	)

	return err
}

func (r *CompetitionRepository) ByID(id uint64) (*app.Competition, error) {
	query := `SELECT * FROM sportmonks_competition where id = $1`
	row := r.connection.QueryRow(query, id)

	return rowToCompetition(row)
}

func (r *CompetitionRepository) Get(q app.CompetitionFilterQuery) ([]app.Competition, error) {
	builder := r.queryBuilder()

	query := builder.Select("sportmonks_competition.*").From("sportmonks_competition")

	if q.IsCup != nil {
		if *q.IsCup {
			query = query.Where(sq.Eq{"is_cup": true})
		} else {
			query = query.Where(sq.Eq{"is_cup": false})
		}
	}

	if len(q.CountryIds) > 0 {
		query = query.Where(sq.Eq{"country_id": q.CountryIds})
	}

	if q.SortBy != nil && *q.SortBy == "id_asc" {
		query = query.OrderBy("id ASC")
	}

	if q.SortBy != nil && *q.SortBy == "id_desc" {
		query = query.OrderBy("id DESC")
	}

	if q.SortBy == nil {
		query = query.OrderBy("id ASC")
	}

	rows, err := query.Query()

	if err != nil {
		return []app.Competition{}, err
	}

	return rowsToCompetition(rows)
}

func (r *CompetitionRepository) IDs() ([]uint64, error) {
	query := `SELECT id FROM sportmonks_competition ORDER BY id ASC`

	rows, err := r.connection.Query(query)

	if err != nil {
		return []uint64{}, err
	}

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

func rowToCompetition(r *sql.Row) (*app.Competition, error) {
	var created int64
	var updated int64

	c := app.Competition{}

	if err := r.Scan(&c.ID, &c.Name, &c.CountryID, &c.IsCup, &created, &updated); err != nil {
		return &c, errors.New(fmt.Sprintf("Competition with ID %d does not exist", c.ID))
	}

	c.CreatedAt = time.Unix(created, 0)
	c.UpdatedAt = time.Unix(updated, 0)

	return &c, nil
}

func rowsToCompetition(rows *sql.Rows) ([]app.Competition, error) {
	defer rows.Close()

	var created int64
	var updated int64
	var comps []app.Competition
	c := app.Competition{}

	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Name, &c.CountryID, &c.IsCup, &created, &updated); err != nil {
			return comps, err
		}

		c.CreatedAt = time.Unix(created, 0)
		c.UpdatedAt = time.Unix(updated, 0)

		comps = append(comps, c)
	}

	return comps, nil
}

func (r *CompetitionRepository) queryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(r.connection)
}

func NewCompetitionRepository(connection *sql.DB, clock clockwork.Clock) *CompetitionRepository {
	return &CompetitionRepository{connection: connection, clock: clock}
}
