package mock

import (
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/mock"
)

type SeasonRepository struct {
	mock.Mock
}

func (m SeasonRepository) Insert(c *model.Season) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m SeasonRepository) Update(c *model.Season) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m SeasonRepository) Id(id int64) (*model.Season, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Season)
	return c, args.Error(1)
}

func (m SeasonRepository) Ids() ([]int64, error) {
	args := m.Called()
	return args.Get(0).([]int64), args.Error(1)
}

func (m SeasonRepository) CurrentSeasonIds() ([]int64, error) {
	args := m.Called()
	return args.Get(0).([]int64), args.Error(1)
}
