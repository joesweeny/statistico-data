package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type SeasonRepository struct {
	mock.Mock
}

func (m SeasonRepository) Insert(c *app.Season) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m SeasonRepository) Update(c *app.Season) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m SeasonRepository) ByID(id int64) (*app.Season, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Season)
	return c, args.Error(1)
}

func (m SeasonRepository) IDs() ([]int64, error) {
	args := m.Called()
	return args.Get(0).([]int64), args.Error(1)
}

func (m SeasonRepository) CurrentSeasonIDs() ([]int64, error) {
	args := m.Called()
	return args.Get(0).([]int64), args.Error(1)
}

type SeasonRequester struct {
	mock.Mock
}

func (s SeasonRequester) Seasons() <-chan *app.Season {
	args := s.Called()
	return args.Get(0).(chan *app.Season)
}
