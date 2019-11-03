package result

import (
	"github.com/statistico/statistico-data/internal/app"
	m "github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandleResult(t *testing.T) {
	teamRepo := new(m.TeamRepository)
	compRepo := new(m.CompetitionRepository)
	roundRepo := new(m.RoundRepository)
	seasonRepo := new(m.SeasonRepository)
	venueRepo := new(m.VenueRepository)
	handler := Handler{
		CompetitionRepo: compRepo,
		RoundRepo:       roundRepo,
		SeasonRepo:      seasonRepo,
		TeamRepo:        teamRepo,
		VenueRepo:       venueRepo,
	}

	t.Run("hydrates new proto result struct", func(t *testing.T) {
		form := "4-3-2-1"
		score := 2
		full := "2-2"
		pos1 := 3
		pos2 := 19
		min := 90
		res := model.Result{}
		res.FixtureID = 92
		res.HomeFormation = &form
		res.AwayFormation = &form
		res.HomeScore = &score
		res.AwayScore = &score
		res.FullTimeScore = &full
		res.HomeLeaguePosition = &pos1
		res.AwayLeaguePosition = &pos2
		res.Minutes = &min

		seasonRepo.On("ByID", uint64(14567)).Return(newSeason(), nil)
		compRepo.On("ByID", uint64(45)).Return(newCompetition(), nil)
		teamRepo.On("ByID", uint64(451)).Return(newTeam(451, "West Ham"), nil)
		teamRepo.On("ByID", uint64(924)).Return(newTeam(924, "Chelsea"), nil)
		venueRepo.On("GetById", uint64(87)).Return(newVenue(), nil)
		roundRepo.On("ByID", uint64(165789)).Return(newRound(), nil)

		proto, err := handler.HandleResult(newFixture(), &res)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(92), proto.Id)
		a.Equal(int64(4), proto.Competition.GetId())
		a.Equal("Premier League", proto.Competition.GetName())
		a.False(proto.Competition.IsCup.GetValue())
		a.Equal(int64(14567), proto.Season.GetId())
		a.Equal("2018-2019", proto.Season.GetName())
		a.True(proto.Season.IsCurrent.GetValue())
		a.Nil(proto.Venue)
		a.Nil(proto.RefereeId)
		a.Equal(int64(1548086929), proto.DateTime)
		a.Equal(int64(451), proto.MatchData.HomeTeam.GetId())
		a.Equal("West Ham", proto.MatchData.HomeTeam.GetName())
		a.Equal(int64(924), proto.MatchData.AwayTeam.GetId())
		a.Equal("Chelsea", proto.MatchData.AwayTeam.GetName())
		a.Nil(proto.MatchData.Stats.Pitch)
		a.Equal("4-3-2-1", proto.MatchData.Stats.HomeFormation.GetValue())
		a.Equal("4-3-2-1", proto.MatchData.Stats.AwayFormation.GetValue())
		a.Equal(int32(2), proto.MatchData.Stats.HomeScore.GetValue())
		a.Equal(int32(2), proto.MatchData.Stats.AwayScore.GetValue())
		a.Nil(proto.MatchData.Stats.HomePenScore)
		a.Nil(proto.MatchData.Stats.AwayPenScore)
		a.Nil(proto.MatchData.Stats.HalfTimeScore)
		a.Nil(proto.MatchData.Stats.ExtraTimeScore)
		a.Equal("2-2", proto.MatchData.Stats.FullTimeScore.GetValue())
		a.Equal(int32(3), proto.MatchData.Stats.HomeLeaguePosition.GetValue())
		a.Equal(int32(19), proto.MatchData.Stats.AwayLeaguePosition.GetValue())
		a.Equal(int32(90), proto.MatchData.Stats.Minutes.GetValue())
		a.Nil(proto.MatchData.Stats.AddedTime)
		a.Nil(proto.MatchData.Stats.ExtraTime)
		a.Nil(proto.MatchData.Stats.InjuryTime)
		a.Equal(int64(165789), proto.Round.GetId())
		a.Equal("18", proto.Round.GetName())
		a.Equal(int64(14567), proto.Round.GetSeasonId())
		a.Equal("2019-01-21T16:08:49Z", proto.Round.GetStartDate())
		a.Equal("2019-01-21T16:08:49Z", proto.Round.GetEndDate())
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
		ID:        87,
		Name:      "London Stadium",
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}

func newFixture() *app.Fixture {
	var roundId = uint64(165789)

	return &app.Fixture{
		ID:         92,
		SeasonID:   14567,
		RoundID:    &roundId,
		HomeTeamID: 451,
		AwayTeamID: 924,
		Date:       time.Unix(1548086929, 0),
		CreatedAt:  time.Unix(1546965200, 0),
		UpdatedAt:  time.Unix(1546965200, 0),
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
