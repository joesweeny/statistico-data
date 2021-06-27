package mock

import (
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type FixtureTeamXGRepository struct {
	mock.Mock
}

func (m *FixtureTeamXGRepository) Insert(f *app.FixtureTeamXG) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *FixtureTeamXGRepository) Update(f *app.FixtureTeamXG) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *FixtureTeamXGRepository) ByID(id uint64) (*app.FixtureTeamXG, error) {
	args := m.Called(id)
	return args.Get(0).(*app.FixtureTeamXG), args.Error(1)
}

func (m *FixtureTeamXGRepository) ByFixtureID(id uint64) (*app.FixtureTeamXG, error) {
	args := m.Called(id)
	return args.Get(0).(*app.FixtureTeamXG), args.Error(1)
}
