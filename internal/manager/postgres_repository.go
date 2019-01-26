package manager

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"time"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("not found")

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

func (p *PostgresManagerRepository) Id (id int) (*model.Manager, error) {
	query := `SELECT * FROM sportmonks_manager where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToManager(row)
}

func rowToManager(r *sql.Row) (*model.Manager, error) {
	var created int64
	var updated int64

	m := model.Manager{}

	err := r.Scan(&m.ID, &m.TeamID, &m.CountryID, &m.FirstName, &m.LastName, &m.Nationality, &m.Image, &created, &updated)

	if err != nil {
		return &m, ErrNotFound
	}

	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}
