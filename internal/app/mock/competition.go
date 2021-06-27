package mock

import (
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type CompetitionRepository struct {
	mock.Mock
}

func (m *CompetitionRepository) Insert(c *app.Competition) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *CompetitionRepository) Update(c *app.Competition) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m *CompetitionRepository) ByID(id uint64) (*app.Competition, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Competition)
	return c, args.Error(1)
}

func (m *CompetitionRepository) Get(q app.CompetitionFilterQuery) ([]app.Competition, error) {
	args := m.Called(q)
	return args.Get(0).([]app.Competition), args.Error(1)
}

func (m *CompetitionRepository) IDs() ([]uint64, error) {
	args := m.Called()
	return args.Get(0).([]uint64), args.Error(1)
}

type CompetitionRequester struct {
	mock.Mock
}

func (c *CompetitionRequester) Competitions() <-chan *app.Competition {
	args := c.Called()
	return args.Get(0).(chan *app.Competition)
}
