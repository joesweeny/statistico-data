package fixture

import (
	"github.com/statistico/statistico-data/internal/app"
	m "github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/season"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestHandleFixture(t *testing.T) {
	teamRepo := new(mockTeamRepository)
	compRepo := new(mockCompetitionRepository)
	roundRepo := new(mockRoundRepository)
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

		seasonRepo.On("Id", int64(14567)).Return(newSeason(), nil)
		compRepo.On("GetById", 45).Return(newCompetition(), nil)
		teamRepo.On("GetById", 451).Return(newTeam(451, "West Ham"), nil)
		teamRepo.On("GetById", 924).Return(newTeam(924, "Chelsea"), nil)
		venueRepo.On("GetById", int64(87)).Return(newVenue(), nil)
		roundRepo.On("GetById", 165789).Return(newRound(), nil)

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

		seasonRepo.On("Id", int64(14567)).Return(newSeason(), nil)
		compRepo.On("GetById", 45).Return(newCompetition(), nil)
		teamRepo.On("GetById", 451).Return(newTeam(451, "West Ham"), nil)
		teamRepo.On("GetById", 924).Return(newTeam(924, "Chelsea"), nil)

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
		teamRepo := new(mockTeamRepository)
		compRepo := new(mockCompetitionRepository)
		roundRepo := new(mockRoundRepository)
		seasonRepo := new(mockSeasonRepository)
		venueRepo := m.VenueRepository{}
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

		seasonRepo.On("Id", int64(14567)).Return(&model.Season{}, season.ErrNotFound)
		compRepo.AssertNotCalled(t, "GetById", 45)
		teamRepo.AssertNotCalled(t, "GetById", 451)
		teamRepo.AssertNotCalled(t, "GetById", 924)
		venueRepo.AssertNotCalled(t, "GetById", 87)
		roundRepo.AssertNotCalled(t, "GetById", 165789)

		proto, err := handler.HandleFixture(fixture)

		if proto != nil {
			t.Fatalf("Test failed, expected nil, got %+v", proto)
		}

		if err == nil {
			t.Fatalf("Test failed, expected %s, got nil", season.ErrNotFound)
		}
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

type mockRoundRepository struct {
	mock.Mock
}

func (m mockRoundRepository) Insert(c *model.Round) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m mockRoundRepository) Update(c *model.Round) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m mockRoundRepository) GetById(id int) (*model.Round, error) {
	args := m.Called(id)
	c := args.Get(0).(*model.Round)
	return c, args.Error(1)
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

func newVenue() *app.Venue {
	return &app.Venue{
		ID:        int64(87),
		Name:      "London Stadium",
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}

func newRound() *model.Round {
	return &model.Round{
		ID:        165789,
		Name:      "18",
		SeasonID:  14567,
		StartDate: time.Unix(1548086929, 0),
		EndDate:   time.Unix(1548086929, 0),
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}
