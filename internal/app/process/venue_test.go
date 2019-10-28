package process_test

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/process"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVenueProcessor_Process(t *testing.T) {
	t.Run("inserts new venue into repository when processing venue command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("GetById", int64(100)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("Insert", whu).Return(nil)
		venueRepo.On("Insert", ars).Return(nil)

		processor.Process("venue", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing venue into repository when processing venue command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(whu, nil)
		venueRepo.On("GetById", int64(100)).Return(ars, nil)
		venueRepo.On("Update", whu).Return(nil)
		venueRepo.On("Update", ars).Return(nil)

		processor.Process("venue", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert venue into repository when processing venue command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("GetById", int64(100)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("Insert", whu).Return(errors.New("error occurred"))
		venueRepo.On("Insert", ars).Return(nil)

		processor.Process("venue", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update venue into repository when processing venue command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(whu, nil)
		venueRepo.On("GetById", int64(100)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("Update", whu).Return(errors.New("error occurred"))
		venueRepo.On("Insert", ars).Return(nil)

		processor.Process("venue", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new venue into repository when processing venue current season command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("GetById", int64(100)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("Insert", whu).Return(nil)
		venueRepo.On("Insert", ars).Return(nil)

		processor.Process("venue:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing venue into repository when processing venue current season command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(whu, nil)
		venueRepo.On("GetById", int64(100)).Return(ars, nil)
		venueRepo.On("Update", whu).Return(nil)
		venueRepo.On("Update", ars).Return(nil)

		processor.Process("venue:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert venue into repository when processing venue current season command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("GetById", int64(100)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("Insert", whu).Return(errors.New("error occurred"))
		venueRepo.On("Insert", ars).Return(nil)

		processor.Process("venue:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update venue into repository when processing venue current season command", func(t *testing.T) {
		t.Helper()
		venueRepo := new(mock.VenueRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.VenueRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewVenueProcessor(venueRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newVenue(44, "London Stadium")
		ars := newVenue(100, "Emirates Stadium")

		venues := make([]*app.Venue, 2)
		venues[0] = whu
		venues[1] = ars

		ch := venueChannel(venues)

		ids := []int64{32}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("VenuesBySeasonIDs", ids).Return(ch)

		venueRepo.On("GetById", int64(44)).Return(whu, nil)
		venueRepo.On("GetById", int64(100)).Return(&app.Venue{}, errors.New("not Found"))
		venueRepo.On("Update", whu).Return(errors.New("error occurred"))
		venueRepo.On("Insert", ars).Return(nil)

		processor.Process("venue:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		venueRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newVenue(id int64, name string) *app.Venue {
	surface := "grass"
	city := "London"
	capacity := 60000

	return &app.Venue{
		ID:       id,
		Name:      name,
		Surface:   &surface,
		City:      &city,
		Capacity:  &capacity,
	}
}

func venueChannel(venues []*app.Venue) chan *app.Venue {
	ch := make(chan *app.Venue, len(venues))

	for _, c := range venues {
		ch <- c
	}

	close(ch)

	return ch
}
