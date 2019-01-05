package country

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/satori/go.uuid"
	"time"
	_ "github.com/lib/pq"
)

type postgresCountryRepository struct {
	Connection *sql.DB
}

func NewPostgresCountryRepository(db *sql.DB) Repository {
	return &postgresCountryRepository{db}
}

func (p *postgresCountryRepository) Insert(c model.Country) error {
	query := `
	INSERT INTO country (id, external_id, name, continent, iso, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.Connection.Exec(
		query,
		c.ID.String(),
		c.ExternalID,
		c.Name,
		c.Continent,
		c.ISO,
		c.CreatedAt.Unix(),
		c.UpdatedAt.Unix(),
	)

	return err
}

func (p *postgresCountryRepository) Update(c model.Country) error {
	_, err := p.GetById(c.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE country set id = $1, external_id = $2, name = $3, continent = $4, iso = $5, updated_at = $6`

	_, err = p.Connection.Exec(
		query,
		c.ID.String(),
		c.ExternalID,
		c.Name,
		c.Continent,
		c.ISO,
		time.Now().Unix(),
	)

	return err
}

func (p *postgresCountryRepository) GetById(u uuid.UUID) (model.Country, error) {
	query := `SELECT * from country where id = $1`
	row := p.Connection.QueryRow(query, u.String())

	return rowToCountry(row)
}

func (p *postgresCountryRepository) GetByExternalId(id int) (model.Country, error) {
	query := `SELECT * from country where external_id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToCountry(row)
}

func rowToCountry(r *sql.Row) (model.Country, error) {
	var id string
	var external int
	var name string
	var continent string
	var iso string
	var created int64
	var updated int64

	c := model.Country{}

	if err := r.Scan(&id, &name, &continent, &iso, &created, &updated, &external); err != nil {
		return c, err
	}

	c.ID = uuid.FromStringOrNil(id)
	c.ExternalID = external
	c.Name = name
	c.Continent = continent
	c.ISO = iso
	c.CreatedAt = time.Unix(created, 0)
	c.UpdatedAt = time.Unix(updated, 0)

	return c, nil
}