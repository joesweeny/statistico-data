package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type VenueRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (r *VenueRepository) Insert(v *app.Venue) error {
	query := `
	INSERT INTO sportmonks_venue (id, name, surface, address, city, capacity, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.connection.Exec(
		query,
		v.ID,
		v.Name,
		v.Surface,
		v.Address,
		v.City,
		v.Capacity,
		r.clock.Now().Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *VenueRepository) Update(v *app.Venue) error {
	_, err := r.GetById(v.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_venue set name = $2, surface = $3, address = $4, city = $5, capacity = $6, updated_at = $7
	WHERE id = $1`

	_, err = r.connection.Exec(
		query,
		v.ID,
		v.Name,
		v.Surface,
		v.Address,
		v.City,
		v.Capacity,
		r.clock.Now().Unix(),
	)

	return err
}

func (r *VenueRepository) GetById(id int64) (*app.Venue, error) {
	query := `SELECT * FROM sportmonks_venue where id = $1`
	row := r.connection.QueryRow(query, id)

	return rowToVenue(row)
}

func rowToVenue(r *sql.Row) (*app.Venue, error) {
	var created int64
	var updated int64

	v := app.Venue{}

	err := r.Scan(&v.ID, &v.Name, &v.Surface, &v.Address, &v.City, &v.Capacity, &created, &updated)

	if err != nil {
		return &v, errors.New(fmt.Sprintf("Venue with ID %d does not exist", v.ID))
	}

	v.CreatedAt = time.Unix(created, 0)
	v.UpdatedAt = time.Unix(updated, 0)

	return &v, nil
}

func NewVenueRepository(connection *sql.DB, clock clockwork.Clock) *VenueRepository {
	return &VenueRepository{connection: connection, clock: clock}
}
