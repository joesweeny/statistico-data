package result

import (
	"github.com/stretchr/testify/mock"
	"github.com/joesweeny/statistico-data/internal/model"
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHandleResult(t *testing.T) {
	teamRepo := new(mockTeamRepository)
	compRepo := new(mockCompetitionRepository)
	seasonRepo := new(mockSeasonRepository)
	venueRepo := new(mockVenueRepository)
	fixtureRepo := new(mockFixtureRepository)
	handler := Handler{
		CompetitionRepo: compRepo,
		SeasonRepo: seasonRepo,
		TeamRepo: teamRepo,
		VenueRepo: venueRepo,
	}

	t.Run("hydrates new proto result struct", func(t *testing.T) {
		form := "4-3-2-1"
		score := 2
		full := "2-2"
		zero := 0
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
		res.Seconds = &zero

		seasonRepo.On("Id", 14567).Return(newSeason(), nil)
		compRepo.On("GetById", 45).Return(newCompetition(), nil)
		teamRepo.On("GetById", 451).Return(newTeam(451, "West Ham"), nil)
		teamRepo.On("GetById", 924).Return(newTeam(924, "Chelsea"), nil)
		venueRepo.On("GetById", 87).Return(newVenue(), nil)

		proto, err := handler.HandleResult(newFixture(), &res)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(92), proto.Id)
		a.Equal(int64(4), proto.Competition.GetId())
		a.Equal("Premier League", proto.Competition.GetName())
		a.False(proto.Competition.GetIsCup())
		a.Equal(int64(14567), proto.Season.GetId())
		a.Equal("2018-2019", proto.Season.GetName())
		a.True(proto.Season.GetIsCurrent())
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
		a.Equal(int32(2), proto.MatchData.Stats.HomeScore)
		a.Equal(int32(2), proto.MatchData.Stats.AwayScore)
		a.Nil(proto.MatchData.Stats.HomePenScore)
		a.Nil(proto.MatchData.Stats.AwayPenScore)
		a.Nil(proto.MatchData.Stats.HalfTimeScore)
		a.Nil(proto.MatchData.Stats.ExtraTimeScore)
		a.Equal("2-2", proto.MatchData.Stats.FullTimeScore.GetValue())
		a.Equal(int32(3), proto.MatchData.Stats.HomeLeaguePosition.GetValue())
		a.Equal(int32(19), proto.MatchData.Stats.AwayLeaguePosition.GetValue())
		a.Equal(int32(90), proto.MatchData.Stats.Minutes.GetValue())
		a.Equal(int32(0), proto.MatchData.Stats.Seconds.GetValue())
		a.Nil(proto.MatchData.Stats.AddedTime)
		a.Nil(proto.MatchData.Stats.ExtraTime)
		a.Nil(proto.MatchData.Stats.InjuryTime)
	})
}

type mockTeamRepository struct {
	mock.Mock
}

func (m mockTeamRepository) Insert(c *model.Team) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockTeamRepository) Update(c *model.Team) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockTeamRepository) GetById(id int) (*model.Team, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Team)
	return c, args.Error(1)
}

type mockCompetitionRepository struct {
	mock.Mock
}

func (m mockCompetitionRepository) Insert(c *model.Competition) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockCompetitionRepository) Update(c *model.Competition) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockCompetitionRepository) GetById(id int) (*model.Competition, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Competition)
	return c, args.Error(1)
}

type mockVenueRepository struct {
	mock.Mock
}

func (m mockVenueRepository) Insert(v *model.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m mockVenueRepository) Update(v *model.Venue) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m mockVenueRepository) GetById(id int) (*model.Venue, error) {
	args := m.Called(id)
	v := args.Get(0).(*model.Venue)
	return v, args.Error(1)
}

type mockSeasonRepository struct {
	mock.Mock
}

func (m mockSeasonRepository) Insert(c *model.Season) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockSeasonRepository) Update(c *model.Season) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockSeasonRepository) Id(id int) (*model.Season, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Season)
	return c, args.Error(1)
}

func (m mockSeasonRepository) Ids() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

func (m mockSeasonRepository) CurrentSeasonIds() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

type mockFixtureRepository struct {
	mock.Mock
}

func newCompetition() *model.Competition {
	return &model.Competition{
		ID:        4,
		Name:      "Premier League",
		CountryID: 462,
		IsCup:     false,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}

func newSeason() *model.Season {
	return &model.Season{
		ID:        14567,
		Name:      "2018-2019",
		LeagueID:  45,
		IsCurrent: true,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}

func newTeam(id int, name string) *model.Team {
	return &model.Team{
		ID:           id,
		Name:         name,
		VenueID:      560,
		NationalTeam: false,
		CreatedAt:    time.Unix(1546965200, 0),
		UpdatedAt:    time.Unix(1546965200, 0),
	}
}

func newVenue() *model.Venue {
	return &model.Venue{
		ID:        87,
		Name:      "London Stadium",
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}

func newFixture() *model.Fixture {
	var roundId = 165789

	return &model.Fixture{
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
