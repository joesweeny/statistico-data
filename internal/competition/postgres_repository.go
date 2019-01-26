package competition

import (
	"database/sql"
	"errors"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresCompetitionRepository struct {
	Connection *sql.DB
}

func (p *PostgresCompetitionRepository) Insert(c *model.Competition) error {
	query := `
	INSERT INTO sportmonks_competition (id, name, country_id, is_cup, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := p.Connection.Exec(
		query,
		c.ID,
		c.Name,
		c.CountryID,
		c.IsCup,
		c.CreatedAt.Unix(),
		c.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresCompetitionRepository) Update(c *model.Competition) error {
	_, err := p.GetById(c.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_competition set name = $2, is_cup = $3, updated_at = $4 where id = $1`

	_, err = p.Connection.Exec(
		query,
		c.ID,
		c.Name,
		c.IsCup,
		c.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresCompetitionRepository) GetById(id int) (*model.Competition, error) {
	query := `SELECT * FROM sportmonks_competition where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToCompetition(row)
}

func rowToCompetition(r *sql.Row) (*model.Competition, error) {
	var created int64
	var updated int64

	c := model.Competition{}

	if err := r.Scan(&c.ID, &c.Name, &c.CountryID, &c.IsCup, &created, &updated); err != nil {
		return &c, ErrNotFound
	}

	c.CreatedAt = time.Unix(created, 0)
	c.UpdatedAt = time.Unix(updated, 0)

	return &c, nil
}
