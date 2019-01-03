package repository

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/country"
	"github.com/joesweeny/statshub/internal/model"
	"errors"
)

type postgresCountryRepository struct {
	Connection *sql.DB
}

func NewPostgresCountryRepository(db *sql.DB) country.Repository {
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
	return errors.New("")
}