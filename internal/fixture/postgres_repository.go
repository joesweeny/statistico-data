package fixture

import (
	"github.com/pkg/errors"
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	"time"
	_ "github.com/lib/pq"
)

var ErrNotFound = errors.New("not found")

type PostgresFixtureRepository struct {
	Connection *sql.DB
}

func (p *PostgresFixtureRepository) Insert(f *model.Fixture) error {
	query := `
	INSERT INTO sportmonks_fixture (id, season_id, round_id, venue_id, home_team_id, away_team_id, referee_id,
	date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.Connection.Exec(
		query,
		f.ID,
		f.SeasonID,
		f.RoundID,
		f.VenueID,
		f.HomeTeamID,
		f.AwayTeamID,
		f.RefereeID,
		f.Date.Unix(),
		f.CreatedAt.Unix(),
		f.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresFixtureRepository) Update(f *model.Fixture) error {
	_, err := p.GetById(f.ID)

	if err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_fixture set season_id = $2, round_id = $3, venue_id = $4, home_team_id = $5, away_team_id = $6,
	referee_id = $7, date = $8, updated_at = $9 where id = $1`

	_, err = p.Connection.Exec(
		query,
		f.ID,
		f.SeasonID,
		f.RoundID,
		f.VenueID,
		f.HomeTeamID,
		f.AwayTeamID,
		f.RefereeID,
		f.Date.Unix(),
		f.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresFixtureRepository) GetById(id int) (*model.Fixture, error) {
	query := `SELECT * FROM sportmonks_fixture where id = $1`
	row := p.Connection.QueryRow(query, id)

	return rowToFixture(row)
}

func rowToFixture(r *sql.Row) (*model.Fixture, error) {
	var id int
	var seasonId int
	var roundId *int
	var venueId *int
	var homeTeam int
	var awayTeam int
	var referee *int
	var date int64
	var created int64
	var updated int64

	f := model.Fixture{}

	err := r.Scan(&id, &seasonId, &roundId, &venueId, &homeTeam, &awayTeam, &referee, &date, &created, &updated)

	if err != nil {
		return &f, ErrNotFound
	}

	f.ID = id
	f.SeasonID = seasonId
	f.RoundID = roundId
	f.VenueID = venueId
	f.HomeTeamID = homeTeam
	f.AwayTeamID = awayTeam
	f.RefereeID = referee
	f.Date = time.Unix(date, 0)
	f.CreatedAt = time.Unix(created, 0)
	f.UpdatedAt = time.Unix(updated, 0)

	return &f, nil
}