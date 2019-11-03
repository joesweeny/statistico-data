package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type PlayerRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (p *PlayerRepository) Insert(m *app.Player) error {
	query := `
	INSERT INTO sportmonks_player (id, country_id, first_name, last_name, birth_place, date_of_birth, position_id, image,
	created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.connection.Exec(
		query,
		m.ID,
		m.CountryId,
		m.FirstName,
		m.LastName,
		m.BirthPlace,
		m.DateOfBirth,
		m.PositionID,
		m.Image,
		p.clock.Now().Unix(),
		p.clock.Now().Unix(),
	)

	return err
}

func (p *PlayerRepository) Update(m *app.Player) error {
	if _, err := p.ByID(m.ID); err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_player set country_id = $2, position_id = $3, image = $4, updated_at = $5 where id = $1`

	_, err := p.connection.Exec(query, m.ID, m.CountryId, m.PositionID, m.Image, p.clock.Now().Unix())

	return err
}

func (p *PlayerRepository) ByID(id uint64) (*app.Player, error) {
	query := `SELECT * FROM sportmonks_player where id = $1`
	row := p.connection.QueryRow(query, id)

	return rowToPlayer(row)
}

func rowToPlayer(r *sql.Row) (*app.Player, error) {
	var created int64
	var updated int64

	var m = app.Player{}

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
		return &m, fmt.Errorf("player with ID %d does not exist", m.ID)
	}

	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}

func NewPlayerRepository(connection *sql.DB, clock clockwork.Clock) *PlayerRepository {
	return &PlayerRepository{connection: connection, clock: clock}
}
