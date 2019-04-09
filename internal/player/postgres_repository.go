package player

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/statistico/statistico-data/internal/model"
	"time"
)

var ErrNotFound = errors.New("not found")

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

func (p *PostgresPlayerRepository) Update(m *model.Player) error {
	if _, err := p.Id(m.ID); err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_player set country_id = $2, position_id = $3, image = $4, updated_at = $5 where id = $1`

	_, err := p.Connection.Exec(query, m.ID, m.CountryId, m.PositionID, m.Image, m.UpdatedAt.Unix())

	return err
}

func (p *PostgresPlayerRepository) Id(id int) (*model.Player, error) {
	query := `SELECT * FROM sportmonks_player where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToPlayer(row)
}

func rowToPlayer(r *sql.Row) (*model.Player, error) {
	var created int64
	var updated int64

	m := model.Player{}

	err := r.Scan(
		&m.ID,
		&m.CountryId,
		&m.FirstName,
		&m.LastName,
		&m.BirthPlace,
		&m.DateOfBirth,
		&m.PositionID,
		&m.Image,
		&created,
		&updated,
	)

	if err != nil {
		return &m, ErrNotFound
	}

	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}
