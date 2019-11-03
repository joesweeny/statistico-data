package postgres

import (
	"database/sql"
	"errors"
	"fmt"
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

func NewCompetitionRepository(connection *sql.DB, clock clockwork.Clock) *CompetitionRepository {
	return &CompetitionRepository{connection: connection, clock: clock}
}
