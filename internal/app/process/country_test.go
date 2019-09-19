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

func TestProcess(t *testing.T) {
	t.Run("inserts new country", func(t *testing.T) {
		repo := new(mock.CountryRepository)
		requester := new(mock.CountryRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCountryProcessor(repo, requester, logger)

		done := make(chan bool)

		eng := newCountry(180, "England")
		ger := newCountry(5, "Germany")

		countries := make([]*app.Country, 2)
		countries[0] = eng
		countries[1] = ger

		ch := countryChannel(countries)

		requester.On("Countries").Return(ch)

		repo.On("GetById", 180).Return(&app.Country{}, errors.New("not Found"))
		repo.On("GetById", 5).Return(&app.Country{}, errors.New("not Found"))
		repo.On("Insert", eng).Return(nil)
		repo.On("Insert", ger).Return(nil)

		processor.Process("country", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing country", func(t *testing.T) {
		repo := new(mock.CountryRepository)
		requester := new(mock.CountryRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCountryProcessor(repo, requester, logger)

		done := make(chan bool)

		eng := newCountry(180, "England")
		ger := newCountry(5, "Germany")

		countries := make([]*app.Country, 2)
		countries[0] = eng
		countries[1] = ger

		ch := countryChannel(countries)

		requester.On("Countries").Return(ch)

		repo.On("GetById", 180).Return(eng, nil)
		repo.On("GetById", 5).Return(ger, nil)
		repo.On("Update", &eng).Return(nil)
		repo.On("Update", &ger).Return(nil)

		processor.Process("country", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error unable to insert country", func(t *testing.T) {
		repo := new(mock.CountryRepository)
		requester := new(mock.CountryRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCountryProcessor(repo, requester, logger)

		done := make(chan bool)

		eng := newCountry(180, "England")
		ger := newCountry(5, "Germany")

		countries := make([]*app.Country, 2)
		countries[0] = eng
		countries[1] = ger

		ch := countryChannel(countries)

		requester.On("Countries").Return(ch)

		repo.On("GetById", 180).Return(&app.Country{}, errors.New("not Found"))
		repo.On("GetById", 5).Return(&app.Country{}, errors.New("not Found"))
		repo.On("Insert", eng).Return(errors.New("error occurred"))
		repo.On("Insert", ger).Return(nil)

		processor.Process("country", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error unable to update country", func(t *testing.T) {
		repo := new(mock.CountryRepository)
		requester := new(mock.CountryRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCountryProcessor(repo, requester, logger)

		done := make(chan bool)

		eng := newCountry(180, "England")
		ger := newCountry(5, "Germany")

		countries := make([]*app.Country, 2)
		countries[0] = eng
		countries[1] = ger

		ch := countryChannel(countries)

		requester.On("Countries").Return(ch)

		repo.On("GetById", 180).Return(eng, nil)
		repo.On("GetById", 5).Return(ger, nil)
		repo.On("Update", &eng).Return(errors.New("error occurred"))
		repo.On("Update", &ger).Return(nil)

		processor.Process("country", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newCountry(id int, name string) *app.Country {
	c := app.Country{
		ID:        id,
		Name:      name,
		Continent: "Europe",
		ISO:       "ENG",
	}

	return &c
}

func countryChannel(countries []*app.Country) chan *app.Country {
	ch := make(chan *app.Country, len(countries))

	for _, c:= range countries {
		ch <- c
	}

	close(ch)

	return ch
}