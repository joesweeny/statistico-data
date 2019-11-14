package handler

import (
	"github.com/sirupsen/logrus/hooks/test"
	m "github.com/statistico/statistico-data/internal/app/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleFixture(t *testing.T) {
	logger, _ := test.NewNullLogger()
	teamRepo := new(m.TeamRepository)
	compRepo := new(m.CompetitionRepository)
	roundRepo := new(m.RoundRepository)
	seasonRepo := new(m.SeasonRepository)
	venueRepo := new(m.VenueRepository)
	handler := FixtureHandler{
		TeamRepo:        teamRepo,
		RoundRepo:       roundRepo,
		VenueRepo:       venueRepo,
		Logger:          logger,
	}

	t.Run("hydrates new proto fixture struct", func(t *testing.T) {
		t.Helper()

		ven := uint64(87)
		ref := uint64(5)
		fixture := newFixture(99)
		fixture.VenueID = &ven
		fixture.RefereeID = &ref

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
		a.Equal(int64(451), proto.HomeTeam.GetId())
		a.Equal("West Ham", proto.HomeTeam.GetName())
		a.Equal(int64(924), proto.AwayTeam.GetId())
		a.Equal("Chelsea", proto.AwayTeam.GetName())
		a.Equal(int64(87), proto.Venue.GetId().GetValue())
		a.Equal("London Stadium", proto.Venue.GetName().GetValue())
		a.Equal(int64(1548086929), proto.GetDateTime().GetUtc())
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
		a.Equal(int64(451), proto.HomeTeam.GetId())
		a.Equal("West Ham", proto.HomeTeam.GetName())
		a.Equal(int64(924), proto.AwayTeam.GetId())
		a.Equal("Chelsea", proto.AwayTeam.GetName())
		a.Nil(proto.Venue.GetId())
		a.Nil(proto.Venue.GetName())
		a.Equal(int64(1548086929), proto.GetDateTime().GetUtc())
		a.Nil(proto.Round)
	})
}
