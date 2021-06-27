package mock

import (
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type RoundRepository struct {
	mock.Mock
}

func (m RoundRepository) Insert(c *app.Round) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m RoundRepository) Update(c *app.Round) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m RoundRepository) ByID(id uint64) (*app.Round, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Round)
	return c, args.Error(1)
}

type RoundRequester struct {
	mock.Mock
}

func (r RoundRequester) RoundsBySeasonIDs(seasonIDs []uint64) <-chan *app.Round {
	args := r.Called(seasonIDs)
	return args.Get(0).(chan *app.Round)
}
