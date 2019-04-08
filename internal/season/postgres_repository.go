package season

import (
	"database/sql"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/pkg/errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresSeasonRepository struct {
	Connection *sql.DB
}

func (p *PostgresSeasonRepository) Insert(s *model.Season) error {
	query := `
	INSERT INTO sportmonks_season (id, name, league_id, is_current, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := p.Connection.Exec(
		query,
		s.ID,
		s.Name,
		s.LeagueID,
		s.IsCurrent,
		s.CreatedAt.Unix(),
		s.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresSeasonRepository) Update(s *model.Season) error {
	_, err := p.Id(s.ID)

	if err != nil {
		return err
	}

	query := `
	 UPDATE sportmonks_season set name = $2, league_id = $3, is_current = $4, updated_at = $5 where id = $1`

	_, err = p.Connection.Exec(
		query,
		s.ID,
		s.Name,
		s.LeagueID,
		s.IsCurrent,
		s.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresSeasonRepository) Id(id int) (*model.Season, error) {
	query := `SELECT * FROM sportmonks_season where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToSeason(row)
}

func (p *PostgresSeasonRepository) Ids() ([]int, error) {
	query := `SELECT id FROM sportmonks_season ORDER BY id ASC`

	rows, err := p.Connection.Query(query)

	if err != nil {
		return []int{}, err
	}

	defer rows.Close()

	var id int
	var ids []int

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (p *PostgresSeasonRepository) CurrentSeasonIds() ([]int, error) {
	query := `SELECT id FROM sportmonks_season where is_current = true ORDER BY id ASC`

	rows, err := p.Connection.Query(query)

	if err != nil {
		return []int{}, err
	}

	defer rows.Close()

	var seasons []int

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return seasons, err
		}

		if err != nil {
			return []int{}, err
		}

		seasons = append(seasons, id)
	}

	return seasons, nil
}

func rowToSeason(r *sql.Row) (*model.Season, error) {
	var created int64
	var updated int64

	s := model.Season{}

	if err := r.Scan(&s.ID, &s.Name, &s.LeagueID, &s.IsCurrent, &created, &updated); err != nil {
		return &s, ErrNotFound
	}

	s.CreatedAt = time.Unix(created, 0)
	s.UpdatedAt = time.Unix(updated, 0)

	return &s, nil
}
