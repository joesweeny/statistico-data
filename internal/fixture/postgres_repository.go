package fixture

import (
	"github.com/pkg/errors"
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
)

var ErrNotFound = errors.New("not found")

type PostgresFixtureRepository struct {
	Connection *sql.DB
}

func (p *PostgresFixtureRepository) Insert(f *model.Fixture) error {

}

func (p *PostgresFixtureRepository) Update(f *model.Fixture) error {

}

func (p *PostgresFixtureRepository) GetById(id int) (*model.Fixture, error) {

}