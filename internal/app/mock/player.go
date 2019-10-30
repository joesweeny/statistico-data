package mock

import (
	"fmt"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type PlayerRepository struct {
	mock.Mock
}

func (m PlayerRepository) Insert(p *app.Player) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m PlayerRepository) Update(p *app.Player) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m PlayerRepository) ByID(id int64) (*app.Player, error) {
	fmt.Println("Hello inside Mock")
	args := m.Called(id)
	return args.Get(0).(*app.Player), args.Error(1)
}

type PlayerRequester struct {
	mock.Mock
}

func (m PlayerRequester) PlayerByID(id int64) *app.Player {
	args := m.Called(id)
	return args.Get(0).(*app.Player)
}
