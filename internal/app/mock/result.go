package mock

import (
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type ResultRepository struct {
	mock.Mock
}

func (m *ResultRepository) Insert(r *app.Result) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *ResultRepository) Update(r *app.Result) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *ResultRepository) ByFixtureID(id uint64) (*app.Result, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Result)
	return c, args.Error(1)
}

type ResultRequester struct {
	mock.Mock
}

func (m *ResultRequester) ResultsBySeasonIDs(id []uint64) <-chan app.Result {
	args := m.Called(id)
	return args.Get(0).(chan app.Result)
}
