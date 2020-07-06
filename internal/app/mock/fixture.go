package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type FixtureRepository struct {
	mock.Mock
}

func (m *FixtureRepository) Insert(c *app.Fixture) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *FixtureRepository) Update(c *app.Fixture) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m *FixtureRepository) ByID(id uint64) (*app.Fixture, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Fixture)
	return c, args.Error(1)
}

func (m *FixtureRepository) ByTeamID(id uint64, query app.FixtureFilterQuery) ([]app.Fixture, error) {
	args := m.Called(id, query)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

func (m *FixtureRepository) Get(q app.FixtureRepositoryQuery) ([]app.Fixture, error) {
	args := m.Called(q)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

func (m *FixtureRepository) GetIDs(q app.FixtureRepositoryQuery) ([]uint64, error) {
	args := m.Called(q)
	return args.Get(0).([]uint64), args.Error(1)
}

type FixtureRequester struct {
	mock.Mock
}

func (m *FixtureRequester) FixturesBySeasonIDs(ids []uint64) <-chan *app.Fixture {
	args := m.Called(ids)
	return args.Get(0).(chan *app.Fixture)
}
