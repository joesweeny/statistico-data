package mock

import (
	"github.com/statistico/statistico-data/internal/statistico"
	"github.com/stretchr/testify/mock"
)

type CountryRepository struct {
	mock.Mock
}

func (m CountryRepository) Insert(c *statistico.Country) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m CountryRepository) Update(c *statistico.Country) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m CountryRepository) GetById(id int) (*statistico.Country, error) {
	args := m.Called(id)
	c := args.Get(0).(*statistico.Country)
	return c, args.Error(1)
}
