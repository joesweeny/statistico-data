package mock

import (
	"github.com/statistico/statistico-data/internal/statistico"
	"github.com/stretchr/testify/mock"
	"time"
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

func Country(id int) *statistico.Country {
	c := statistico.Country{
		ID:        id,
		Name:      "England",
		Continent: "Europe",
		ISO:       "ENG",
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}

	return &c
}