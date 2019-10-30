package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/lib/pq"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
	"time"
)

type SquadRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (r *SquadRepository) Insert(s *app.Squad) error {
	query := `
	INSERT INTO sportmonks_squad (season_id, team_id, player_ids, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.connection.Exec(
		query,
		s.SeasonID,
		s.TeamID,
		pq.Array(s.PlayerIDs),
		r.clock.Now().Unix(),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *SquadRepository) Update(s *app.Squad) error {
	if _, err := r.BySeasonAndTeam(s.SeasonID, s.TeamID); err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_squad set player_ids = $3, updated_at = $4 where season_id = $1 and team_id = $2`

	_, err := r.connection.Exec(
		query,
		s.SeasonID,
		s.TeamID,
		pq.Array(s.PlayerIDs),
		r.clock.Now().Unix(),
	)

	return err
}

func (r *SquadRepository) BySeasonAndTeam(seasonId, teamId int64) (*app.Squad, error) {
	query := `SELECT * FROM sportmonks_squad where season_id = $1 AND team_id = $2`

	var s = app.Squad{}
	var players []string
	var created int64
	var updated int64

	err := r.connection.QueryRow(query, seasonId, teamId).Scan(&s.SeasonID, &s.TeamID, pq.Array(&players), &created, &updated)

	if err != nil {
		return &s, fmt.Errorf("squad with Team ID %d and Season ID %s does not exist", teamId, seasonId)
	}

	for _, i := range players {
		text, _ := strconv.Atoi(i)
		s.PlayerIDs = append(s.PlayerIDs, int64(text))
	}

	s.CreatedAt = time.Unix(created, 0)
	s.UpdatedAt = time.Unix(updated, 0)

	return &s, err
}

func (r *SquadRepository) All() ([]app.Squad, error) {
	query := `SELECT * FROM sportmonks_squad order by season_id ASC, team_id ASC`

	var squads []app.Squad

	rows, err := r.connection.Query(query)

	if err != nil {
		return squads, err
	}

	return parseRows(rows, squads)
}

func (r *SquadRepository) CurrentSeason() ([]app.Squad, error) {
	query := `SELECT * FROM sportmonks_squad WHERE season_id in (SELECT id from sportmonks_season WHERE is_current = true)
 	order by season_id ASC, team_id ASC`

	var squads []app.Squad

	rows, err := r.connection.Query(query)

	if err != nil {
		return squads, err
	}

	return parseRows(rows, squads)
}

func parseRows(r *sql.Rows, m []app.Squad) ([]app.Squad, error) {
	for r.Next() {
		var players []string
		var created int64
		var updated int64
		var squad app.Squad

		if err := r.Scan(&squad.SeasonID, &squad.TeamID, pq.Array(&players), &created, &updated); err != nil {
			return m, err
		}

		for _, i := range players {
			text, _ := strconv.Atoi(i)
			squad.PlayerIDs = append(squad.PlayerIDs, text)
		}

		squad.CreatedAt = time.Unix(created, 0)
		squad.UpdatedAt = time.Unix(updated, 0)

		m = append(m, squad)
	}

	return m, nil
}
