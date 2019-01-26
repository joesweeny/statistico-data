package round

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresRoundRepository struct {
	Connection *sql.DB
}

func (p *PostgresRoundRepository) Insert(r *model.Round) error {
	query := `
	INSERT INTO sportmonks_round (id, name, season_id, start_date, end_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.Connection.Exec(
		query,
		r.ID,
		r.Name,
		r.SeasonID,
		r.StartDate.Unix(),
		r.EndDate.Unix(),
		r.CreatedAt.Unix(),
		r.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresRoundRepository) GetById(id int) (*model.Round, error) {
	query := `SELECT * FROM sportmonks_round where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToRound(row)
}

func (p *PostgresRoundRepository) Update(r *model.Round) error {
	_, err := p.GetById(r.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_round set name = $2, season_id = $3, start_date = $4, end_date = $5, updated_at = $6
	where id = $1`

	_, err = p.Connection.Exec(
		query,
		r.ID,
		r.Name,
		r.SeasonID,
		r.StartDate.Unix(),
		r.EndDate.Unix(),
		r.UpdatedAt.Unix(),
	)

	return err
}

func rowToRound(r *sql.Row) (*model.Round, error) {
	var start int64
	var end int64
	var created int64
	var updated int64

	m := model.Round{}

	err := r.Scan(&m.ID, &m.Name, &m.SeasonID, &start, &end, &created, &updated)

	if err != nil {
		return &m, ErrNotFound
	}

	m.StartDate = time.Unix(start, 0)
	m.EndDate = time.Unix(end, 0)
	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}
