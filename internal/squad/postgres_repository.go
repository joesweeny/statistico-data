package squad

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresSquadRepository struct {
	Connection *sql.DB
}

func (p *PostgresSquadRepository) Insert(m *model.Squad) error {
	query := `
	INSERT INTO sportmonks_squad (season_id, team_id, player_ids, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := p.Connection.Exec(query, m.SeasonID, m.TeamID, pq.Array(m.PlayerIDs), m.CreatedAt.Unix(), m.UpdatedAt.Unix())

	return err
}

func (p *PostgresSquadRepository) Update(m *model.Squad) error {
	if _, err := p.BySeasonAndTeam(m.SeasonID, m.TeamID); err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_squad set player_ids = $3, updated_at = $4 where season_id = $1 and team_id = $2`

	_, err := p.Connection.Exec(query, m.SeasonID, m.TeamID, pq.Array(m.PlayerIDs), m.UpdatedAt.Unix())

	return err
}

func (p *PostgresSquadRepository) BySeasonAndTeam(seasonId, teamId int) (*model.Squad, error) {
	query := `SELECT * FROM sportmonks_squad where season_id = $1 AND team_id = $2`

	var players []string
	var m model.Squad
	var created int64
	var updated int64

	err := p.Connection.QueryRow(query, seasonId, teamId).Scan(&m.SeasonID, &m.TeamID, pq.Array(&players), &created, &updated)

	if err != nil {
		return &model.Squad{}, ErrNotFound
	}

	for _, i := range players {
		text, _ := strconv.Atoi(i)
		m.PlayerIDs = append(m.PlayerIDs, text)
	}

	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, err
}
