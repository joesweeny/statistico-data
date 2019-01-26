package manager

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
)

type PostgresManagerRepository struct {
	Connection *sql.DB
}

func (p *PostgresManagerRepository) Insert(m *model.Manager) error {
	query := `
	INSERT INTO sportmonks_manager (id, team_id, country_id, first_name, last_name, nationality, image, created_at, 
	updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := p.Connection.Exec(
		query,
		m.ID,
		m.TeamID,
		m.CountryID,
		m.FirstName,
		m.LastName,
		m.Nationality,
		m.Image,
		m.CreatedAt.Unix(),
		m.UpdatedAt.Unix(),
	)

	return err
}
