package country

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type PostgresCountryRepository struct {
	Connection *sql.DB
}

var ErrNotFound = errors.New("not found")

func (p *PostgresCountryRepository) Insert(c *model.Country) error {
	query := `
	INSERT INTO sportmonks_country (id, name, continent, iso, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := p.Connection.Exec(
		query,
		c.ID,
		c.Name,
		c.Continent,
		c.ISO,
		c.CreatedAt.Unix(),
		c.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresCountryRepository) Update(c *model.Country) error {
	_, err := p.GetById(c.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_country set name = $2, continent = $3, iso = $4, updated_at = $5 where id = $1`

	_, err = p.Connection.Exec(
		query,
		c.ID,
		c.Name,
		c.Continent,
		c.ISO,
		c.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresCountryRepository) GetById(id int) (*model.Country, error) {
	query := `SELECT * from sportmonks_country where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToCountry(row)
}

func rowToCountry(r *sql.Row) (*model.Country, error) {
	var created int64
	var updated int64

	c := model.Country{}

	if err := r.Scan(&c.ID, &c.Name, &c.Continent, &c.ISO, &created, &updated); err != nil {
		return &c, ErrNotFound
	}

	c.CreatedAt = time.Unix(created, 0)
	c.UpdatedAt = time.Unix(updated, 0)

	return &c, nil
}
