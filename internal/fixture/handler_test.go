package fixture

import (
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/season"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestHandleFixture(t *testing.T) {
	teamRepo := new(mockTeamRepository)
	compRepo := new(mockCompetitionRepository)
	seasonRepo := new(mockSeasonRepository)
	venueRepo := new(mockVenueRepository)
	handler := Handler{
		TeamRepo:        teamRepo,
		CompetitionRepo: compRepo,
		SeasonRepo:      seasonRepo,
		VenueRepo:       venueRepo,
	}

	t.Run("hydrates new proto fixture struct", func(t *testing.T) {
		ven := 87
		ref := 5
		fixture := newFixture(99)
		fixture.VenueID = &ven
		fixture.RefereeID = &ref

		seasonRepo.On("Id", 14567).Return(newSeason(), nil)
		compRepo.On("GetById", 45).Return(newCompetition(), nil)
		teamRepo.On("GetById", 451).Return(newTeam(451, "West Ham"), nil)
		teamRepo.On("GetById", 924).Return(newTeam(924, "Chelsea"), nil)
		venueRepo.On("GetById", 87).Return(newVenue(), nil)

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
	})

	t.Run("can handle nullable fields", func(t *testing.T) {
		fixture := newFixture(99)

		seasonRepo.On("Id", 14567).Return(newSeason(), nil)
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
	})

	t.Run("error is returned if season not found", func(t *testing.T) {
		teamRepo := new(mockTeamRepository)
		compRepo := new(mockCompetitionRepository)
		seasonRepo := new(mockSeasonRepository)
		venueRepo := new(mockVenueRepository)
		handler := Handler{
			TeamRepo:        teamRepo,
			CompetitionRepo: compRepo,
			SeasonRepo:      seasonRepo,
			VenueRepo:       venueRepo,
		}

		ven := 87
		fixture := newFixture(99)
		fixture.VenueID = &ven

		seasonRepo.On("Id", 14567).Return(&model.Season{}, season.ErrNotFound)
		compRepo.AssertNotCalled(t, "GetById", 45)
		teamRepo.AssertNotCalled(t, "GetById", 451)
		teamRepo.AssertNotCalled(t, "GetById", 924)
		venueRepo.AssertNotCalled(t, "GetById", 87)

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
