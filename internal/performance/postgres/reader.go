package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-data/internal/performance"
)

type StatReader struct {
	connection *sql.DB
}

func (s *StatReader) GetTeams(f *performance.StatFilter) ([]*performance.Team, error) {
	builder := s.queryBuilder()

	query := builder.Select("team_id, team_name")

	rows, err := BuildQuery(query, f).Query()

	if err != nil {
		return nil, nil
	}

	var teams []*performance.Team
	var t performance.Team

	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Name)

		if err != nil {
			return nil, err
		}

		teams = append(teams, &t)
	}

	return teams, nil
}

func (s *StatReader) queryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(s.connection)
}