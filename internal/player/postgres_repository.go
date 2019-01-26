package player

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
)

type PostgresPlayerRepository struct {
	Connection *sql.DB
}

func (p *PostgresPlayerRepository) Insert(m *model.Player) error {
	query := `
	INSERT INTO sportmonks_player (id, country_id, first_name, last_name, birth_place, date_of_birth, position_id, image,
	created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.Connection.Exec(
		query,
		m.ID,
		m.CountryId,
		m.FirstName,
		m.LastName,
		m.BirthPlace,
		m.DateOfBirth,
		m.PositionID,
		m.Image,
		m.CreatedAt.Unix(),
		m.UpdatedAt.Unix(),
	)

	return err
}
