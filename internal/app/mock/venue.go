package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type VenueRepository struct {
	mock.Mock
}

func (m VenueRepository) Insert(v *app.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m VenueRepository) Update(v *app.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m VenueRepository) GetById(id uint64) (*app.Venue, error) {
	args := m.Called(id)
	v := args.Get(0).(*app.Venue)
	return v, args.Error(1)
}

type VenueRequester struct {
	mock.Mock
}

func (v VenueRequester) VenuesBySeasonIDs(seasonIDs []uint64) <-chan *app.Venue {
	args := v.Called(seasonIDs)
	return args.Get(0).(chan *app.Venue)
}
