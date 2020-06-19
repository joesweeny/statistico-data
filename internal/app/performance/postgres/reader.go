package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-data/internal/app/performance"
)

type StatReader struct {
	connection *sql.DB
}

func (s *StatReader) TeamsMatchingFilter(f *performance.StatFilter) ([]*performance.Team, error) {
	rows, err := buildTeamsQuery(s.queryBuilder(), f).Query()

	if err != nil {
		return nil, nil
	}

	var teams []*performance.Team

	for rows.Next() {
		var t performance.Team

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

func NewStatReader(connection *sql.DB) *StatReader {
	return &StatReader{connection: connection}
}
