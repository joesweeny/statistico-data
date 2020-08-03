package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type SeasonRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (r *SeasonRepository) Insert(s *app.Season) error {
	query := `
	INSERT INTO sportmonks_season (id, name, league_id, is_current, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.connection.Exec(
		query,
		s.ID,
		s.Name,
		s.CompetitionID,
		s.IsCurrent,
		r.clock.Now().Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *SeasonRepository) Update(s *app.Season) error {
	_, err := r.ByID(s.ID)

	if err != nil {
		return err
	}

	query := `
	 UPDATE sportmonks_season set name = $2, league_id = $3, is_current = $4, updated_at = $5 where id = $1`

	_, err = r.connection.Exec(
		query,
		s.ID,
		s.Name,
		s.CompetitionID,
		s.IsCurrent,
		r.clock.Now().Unix(),
	)

	return err
}

func (r *SeasonRepository) ByID(id uint64) (*app.Season, error) {
	query := `SELECT * FROM sportmonks_season where id = $1`
	row := r.connection.QueryRow(query, id)

	return rowToSeason(row)
}

func (r *SeasonRepository) IDs() ([]uint64, error) {
	query := `SELECT id FROM sportmonks_season ORDER BY id ASC`

	rows, err := r.connection.Query(query)

	if err != nil {
		return []uint64{}, err
	}

	defer rows.Close()

	var id uint64
	var ids []uint64

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *SeasonRepository) CurrentSeasonIDs() ([]uint64, error) {
	query := `SELECT id FROM sportmonks_season where is_current = true ORDER BY id ASC`

	rows, err := r.connection.Query(query)

	if err != nil {
		return []uint64{}, err
	}

	defer rows.Close()

	var seasons []uint64

	for rows.Next() {
		var id uint64

		if err := rows.Scan(&id); err != nil {
			return seasons, err
		}

		seasons = append(seasons, id)
	}

	return seasons, nil
}

func (r *SeasonRepository) ByCompetitionId(id uint64, sort string) ([]app.Season, error) {
	builder := r.queryBuilder()

	query := builder.Select("*").
		From("sportmonks_season").
		Where(sq.Eq{"league_id": id})

	if sort == "name_asc" {
		query = query.OrderBy("name ASC")
	}

	if sort == "name_desc" {
		query = query.OrderBy("name DESC")
	}

	rows, err := query.Query()

	if err != nil {
		return []app.Season{}, err
	}

	var created int64
	var updated int64
	var seasons []app.Season
	var season app.Season

	for rows.Next() {
		err := rows.Scan(
			&season.ID,
			&season.Name,
			&season.CompetitionID,
			&season.IsCurrent,
			&created, &updated,
		)

		if err != nil {
			return seasons, err
		}

		season.CreatedAt = time.Unix(created, 0)
		season.UpdatedAt = time.Unix(updated, 0)

		seasons = append(seasons, season)
	}

	return seasons, nil
}

func (r *SeasonRepository) ByTeamId(id uint64, sort string) ([]app.Season, error) {
	//
}

func rowToSeason(r *sql.Row) (*app.Season, error) {
	var created int64
	var updated int64

	var s = app.Season{}

	if err := r.Scan(&s.ID, &s.Name, &s.CompetitionID, &s.IsCurrent, &created, &updated); err != nil {
		return &s, err
	}

	s.CreatedAt = time.Unix(created, 0)
	s.UpdatedAt = time.Unix(updated, 0)

	return &s, nil
}

func (r *SeasonRepository) queryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(r.connection)
}

func NewSeasonRepository(connection *sql.DB, clock clockwork.Clock) *SeasonRepository {
	return &SeasonRepository{connection: connection, clock: clock}
}
