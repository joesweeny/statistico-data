package country

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	"time"
	_ "github.com/lib/pq"
)

type PostgresCountryRepository struct {
	Connection *sql.DB
}

func (p *PostgresCountryRepository) Insert(c model.Country) error {
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

func (p *PostgresCountryRepository) Update(c model.Country) error {
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

func (p *PostgresCountryRepository) GetById(id int) (model.Country, error) {
	query := `SELECT * from sportmonks_country where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToCountry(row)
}

func rowToCountry(r *sql.Row) (model.Country, error) {
	var id int
	var name string
	var continent string
	var iso string
	var created int64
	var updated int64

	c := model.Country{}

	if err := r.Scan(&id, &name, &continent, &iso, &created, &updated); err != nil {
		return c, err
	}

	c.ID = id
	c.Name = name
	c.Continent = continent
	c.ISO = iso
	c.CreatedAt = time.Unix(created, 0)
	c.UpdatedAt = time.Unix(updated, 0)

	return c, nil
}