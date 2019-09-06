package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type CountryRepository struct {
	mock.Mock
}

func (m CountryRepository) Insert(c *app.Country) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m CountryRepository) Update(c *app.Country) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m CountryRepository) GetById(id int) (*app.Country, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.Country)
	return c, args.Error(1)
}
