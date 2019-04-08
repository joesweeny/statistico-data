package venue

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/statistico/statistico-data/internal/model"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresVenueRepository struct {
	Connection *sql.DB
}

func (p *PostgresVenueRepository) Insert(v *model.Venue) error {
	query := `
	INSERT INTO sportmonks_venue (id, name, surface, address, city, capacity, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := p.Connection.Exec(
		query,
		v.ID,
		v.Name,
		v.Surface,
		v.Address,
		v.City,
		v.Capacity,
		v.CreatedAt.Unix(),
		v.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresVenueRepository) Update(v *model.Venue) error {
	_, err := p.GetById(v.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_venue set name = $2, surface = $3, address = $4, city = $5, capacity = $6, updated_at = $7
	WHERE id = $1`

	_, err = p.Connection.Exec(
		query,
		v.ID,
		v.Name,
		v.Surface,
		v.Address,
		v.City,
		v.Capacity,
		v.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresVenueRepository) GetById(id int) (*model.Venue, error) {
	query := `SELECT * FROM sportmonks_venue where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToVenue(row)
}

func rowToVenue(r *sql.Row) (*model.Venue, error) {
	var created int64
	var updated int64

	v := model.Venue{}

	err := r.Scan(&v.ID, &v.Name, &v.Surface, &v.Address, &v.City, &v.Capacity, &created, &updated)

	if err != nil {
		return &v, ErrNotFound
	}

	v.CreatedAt = time.Unix(created, 0)
	v.UpdatedAt = time.Unix(updated, 0)

	return &v, nil

}
