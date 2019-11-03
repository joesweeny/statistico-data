package fixture

import (
	"errors"
	"github.com/statistico/statistico-data/internal/app"
	m "github.com/statistico/statistico-data/internal/app/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestHandleFixture(t *testing.T) {
	teamRepo := new(m.TeamRepository)
	compRepo := new(m.CompetitionRepository)
	roundRepo := new(m.RoundRepository)
	seasonRepo := new(m.SeasonRepository)
	venueRepo := new(m.VenueRepository)
	handler := Handler{
		TeamRepo:        teamRepo,
		CompetitionRepo: compRepo,
		RoundRepo:       roundRepo,
		SeasonRepo:      seasonRepo,
		VenueRepo:       venueRepo,
		Logger:          log.New(ioutil.Discard, "Error: ", 0),
	}

	t.Run("hydrates new proto fixture struct", func(t *testing.T) {
		t.Helper()

		ven := 87
		ref := 5
		fixture := newFixture(99)
		fixture.VenueID = &ven
		fixture.RefereeID = &ref

		seasonRepo.On("ByID", uint64(14567)).Return(newSeason(), nil)
		compRepo.On("ByID", uint64(45)).Return(newCompetition(), nil)
		teamRepo.On("ByID", uint64(451)).Return(newTeam(451, "West Ham"), nil)
		teamRepo.On("ByID", uint64(924)).Return(newTeam(924, "Chelsea"), nil)
		venueRepo.On("GetById", uint64(87)).Return(newVenue(), nil)
		roundRepo.On("ByID", uint64(165789)).Return(newRound(), nil)

		proto, err := handler.HandleFixture(fixture)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(99), proto.Id)
		a.Equal(int64(4), proto.Competition.GetId())
		a.Equal("Premier League", proto.Competition.GetName())
		a.False(proto.Competition.IsCup.GetValue())
		a.Equal(int64(14567), proto.Season.GetId())
		a.Equal("2018-2019", proto.Season.GetName())
		a.True(proto.Season.IsCurrent.GetValue())
		a.Equal(int64(451), proto.HomeTeam.GetId())
		a.Equal("West Ham", proto.HomeTeam.GetName())
		a.Equal(int64(924), proto.AwayTeam.GetId())
		a.Equal("Chelsea", proto.AwayTeam.GetName())
		a.Equal(int64(87), proto.Venue.GetId().GetValue())
		a.Equal("London Stadium", proto.Venue.GetName().GetValue())
		a.Equal(int64(5), proto.RefereeId.GetValue())
		a.Equal(int64(1548086929), proto.GetDateTime())
		a.Equal(int64(165789), proto.Round.GetId())
		a.Equal("18", proto.Round.GetName())
		a.Equal(int64(14567), proto.Round.GetSeasonId())
		a.Equal("2019-01-21T16:08:49Z", proto.Round.GetStartDate())
		a.Equal("2019-01-21T16:08:49Z", proto.Round.GetEndDate())
	})

	t.Run("can handle nullable fields", func(t *testing.T) {
		fixture := newFixture(99)
		fixture.RoundID = nil

		seasonRepo.On("ByID", uint64(14567)).Return(newSeason(), nil)
		compRepo.On("ByID", uint64(45)).Return(newCompetition(), nil)
		teamRepo.On("ByID", uint64(451)).Return(newTeam(451, "West Ham"), nil)
		teamRepo.On("ByID", uint64(924)).Return(newTeam(924, "Chelsea"), nil)

		proto, err := handler.HandleFixture(fixture)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(99), proto.Id)
		a.Equal(int64(4), proto.Competition.GetId())
		a.Equal("Premier League", proto.Competition.GetName())
		a.Equal(int64(14567), proto.Season.GetId())
		a.Equal("2018-2019", proto.Season.GetName())
		a.Equal(int64(451), proto.HomeTeam.GetId())
		a.Equal("West Ham", proto.HomeTeam.GetName())
		a.Equal(int64(924), proto.AwayTeam.GetId())
		a.Equal("Chelsea", proto.AwayTeam.GetName())
		a.Nil(proto.Venue.GetId())
		a.Nil(proto.Venue.GetName())
		a.Nil(proto.RefereeId)
		a.Equal(int64(1548086929), proto.GetDateTime())
		a.Nil(proto.Round)
	})

	t.Run("error is returned if season not found", func(t *testing.T) {
		teamRepo := new(m.TeamRepository)
		compRepo := new(m.CompetitionRepository)
		roundRepo := new(m.RoundRepository)
		seasonRepo := new(m.SeasonRepository)
		venueRepo := new(m.VenueRepository)
		handler := Handler{
			TeamRepo:        teamRepo,
			CompetitionRepo: compRepo,
			RoundRepo:       roundRepo,
			SeasonRepo:      seasonRepo,
			VenueRepo:       venueRepo,
			Logger:          log.New(ioutil.Discard, "Error: ", 0),
		}

		ven := 87
		fixture := newFixture(99)
		fixture.VenueID = &ven

		seasonRepo.On("ByID", uint64(14567)).Return(&app.Season{}, errors.New("not found"))
		compRepo.AssertNotCalled(t, "GetById", 45)
		teamRepo.AssertNotCalled(t, "ByID", int64(451))
		teamRepo.AssertNotCalled(t, "ByID", int64(924))
		venueRepo.AssertNotCalled(t, "GetById", 87)
		roundRepo.AssertNotCalled(t, "ByID", int64(165789))

		proto, err := handler.HandleFixture(fixture)

		if proto != nil {
			t.Fatalf("Test failed, expected nil, got %+v", proto)
		}

		if err == nil {
			t.Fatalf("Test failed, expected %s, got nil", errors.New("not found"))
		}
	})
}

func newCompetition() *app.Competition {
	return &app.Competition{
		ID:        4,
		Name:      "Premier League",
		CountryID: 462,
		IsCup:     false,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}

func newSeason() *app.Season {
	return &app.Season{
		ID:            uint64(14567),
		Name:          "2018-2019",
		CompetitionID: uint64(45),
		IsCurrent:     true,
		CreatedAt:     time.Unix(1546965200, 0),
		UpdatedAt:     time.Unix(1546965200, 0),
	}
}

func newTeam(id int, name string) *app.Team {
	return &app.Team{
		ID:           uint64(id),
		Name:         name,
		VenueID:      uint64(560),
		NationalTeam: false,
		CreatedAt:    time.Unix(1546965200, 0),
		UpdatedAt:    time.Unix(1546965200, 0),
	}
}

func newVenue() *app.Venue {
	return &app.Venue{
		ID:        uint64(87),
		Name:      "London Stadium",
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}

func newRound() *app.Round {
	return &app.Round{
		ID:        165789,
		Name:      "18",
		SeasonID:  14567,
		StartDate: time.Unix(1548086929, 0),
		EndDate:   time.Unix(1548086929, 0),
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}
