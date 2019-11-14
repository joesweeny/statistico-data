package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
	"time"
)

type FixtureRepository struct {
	mock.Mock
}

func (m FixtureRepository) Insert(c *app.Fixture) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m FixtureRepository) Update(c *app.Fixture) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m FixtureRepository) ByID(id uint64) (*app.Fixture, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Fixture)
	return c, args.Error(1)
}

func (m FixtureRepository) IDs() ([]uint64, error) {
	args := m.Called()
	return args.Get(0).([]uint64), args.Error(1)
}

func (m FixtureRepository) IDsBetween(from, to time.Time) ([]uint64, error) {
	args := m.Called(from, to)
	return args.Get(0).([]uint64), args.Error(1)
}

func (m FixtureRepository) Between(from, to time.Time) ([]app.Fixture, error) {
	args := m.Called(from, to)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

func (m FixtureRepository) ByTeamID(id uint64, limit int32, before time.Time) ([]app.Fixture, error) {
	args := m.Called(id, limit, before)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

func (m FixtureRepository) BySeasonID(id uint64) ([]app.Fixture, error) {
	args := m.Called(id)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

func (m FixtureRepository) BySeasonIDBefore(id uint64, before time.Time) ([]app.Fixture, error) {
	args := m.Called(id)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

func (m FixtureRepository) BySeasonIDBetween(id uint64, after, before time.Time) ([]app.Fixture, error) {
	args := m.Called(id, after, before)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

func (m FixtureRepository) ByHomeAndAwayTeam(homeTeamId, awayTeamId uint64, limit uint32, before time.Time) ([]app.Fixture, error) {
	args := m.Called(homeTeamId, awayTeamId, limit, before)
	return args.Get(0).([]app.Fixture), args.Error(1)
}

type FixtureRequester struct {
	mock.Mock
}

func (m FixtureRequester) FixturesBySeasonIDs(ids []uint64) <-chan *app.Fixture {
	args := m.Called(ids)
	return args.Get(0).(chan *app.Fixture)
}
