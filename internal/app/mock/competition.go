package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type CompetitionRepository struct {
	mock.Mock
}

func (m CompetitionRepository) Insert(c *app.Competition) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m CompetitionRepository) Update(c *app.Competition) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m CompetitionRepository) ByID(id int64) (*app.Competition, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Competition)
	return c, args.Error(1)
}

type CompetitionRequester struct {
	mock.Mock
}

func (c CompetitionRequester) Competitions() <-chan *app.Competition {
	args := c.Called()
	return args.Get(0).(chan *app.Competition)
}
