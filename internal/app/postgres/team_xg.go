package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type FixtureTeamXGRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (r *FixtureTeamXGRepository) Insert(f *app.FixtureTeamXG) error {
	query := `
	INSERT INTO understat_fixture_team_xg (id, sportmonks_fixture_id, home, away, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.connection.Exec(
		query,
		f.ID,
		f.FixtureID,
		f.Home,
		f.Away,
		r.clock.Now().Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *FixtureTeamXGRepository) Update(f *app.FixtureTeamXG) error {
	query := `SELECT EXISTS(SELECT 1 from understat_fixture_team_xg where id = $1)`

	var exists bool

	if err := r.connection.QueryRow(query, f.ID).Scan(&exists); err != nil || exists == false {
		return fmt.Errorf("fixture team XG with ID %d does not exist", f.ID)
	}

	query = `UPDATE understat_fixture_team_xg set home = $2, away = $3, updated_at = $4 where id = $1`

	_, err := r.connection.Exec(query, f.ID, f.Home, f.Away, r.clock.Now().Unix())

	return err
}

func (r *FixtureTeamXGRepository) ByID(id uint64) (*app.FixtureTeamXG, error) {
	query := `SELECT * FROM understat_fixture_team_xg where id = $1`

	row := r.connection.QueryRow(query, id)

	return rowToFixtureTeamXG(row, id)
}

func (r *FixtureTeamXGRepository) ByFixtureID(id uint64) (*app.FixtureTeamXG, error) {
	query := `SELECT * FROM understat_fixture_team_xg where sportmonks_fixture_id = $1`

	row := r.connection.QueryRow(query, id)

	return rowToFixtureTeamXG(row, id)
}

func rowToFixtureTeamXG(r *sql.Row, id uint64) (*app.FixtureTeamXG, error) {
	var created int64
	var updated int64

	var f app.FixtureTeamXG

	if err := r.Scan(&f.ID, &f.FixtureID, &f.Home, &f.Away, &created, &updated); err != nil {
		return &f, fmt.Errorf("error fetching fixture team xg from database. ID %d. Error %s", id, err.Error())
	}

	f.CreatedAt = time.Unix(created, 0)
	f.UpdatedAt = time.Unix(updated, 0)

	return &f, nil
}

func NewFixtureTeamXGRepository(connection *sql.DB, clock clockwork.Clock) *FixtureTeamXGRepository {
	return &FixtureTeamXGRepository{connection: connection, clock: clock}
}
