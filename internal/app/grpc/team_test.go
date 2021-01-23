package grpc_test

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	e "github.com/statistico/statistico-data/internal/app/errors"
	"github.com/statistico/statistico-data/internal/app/grpc"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTeamService_GetTeamById(t *testing.T) {
	t.Run("returns a proto Team struct", func(t *testing.T) {
		t.Helper()

		request := statistico.TeamRequest{TeamId: 1}

		repo := new(mock.TeamRepository)
		logger, _ := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)

		code := "WHU"
		founded := 1895
		logo := "https://logo.com"

		team := app.Team{
			ID:           1,
			Name:         "West Ham United",
			ShortCode:    &code,
			CountryID:    8,
			VenueID:      214,
			NationalTeam: false,
			Founded:      &founded,
			Logo:         &logo,
		}

		repo.On("ByID", uint64(1)).Return(&team, nil)

		response, err := service.GetTeamByID(context.Background(), &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(uint64(1), response.Id)
		a.Equal("West Ham United", response.Name)
		a.Equal("WHU", response.ShortCode.Value)
		a.Equal(uint64(8), response.CountryId)
		a.Equal(uint64(214), response.VenueId)
		a.Equal(false, response.IsNationalTeam.Value)
		a.Equal(uint64(1895), response.Founded.Value)
		a.Equal("https://logo.com", response.Logo.Value)
	})

	t.Run("returns not found if error returned by repository", func(t *testing.T) {
		t.Helper()

		request := statistico.TeamRequest{TeamId: 1}

		repo := new(mock.TeamRepository)
		logger, _ := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)

		repo.On("ByID", uint64(1)).Return(&app.Team{}, e.ErrorNotFound)

		_, err := service.GetTeamByID(context.Background(), &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = NotFound desc = team with ID 1 does not exist", err.Error())
	})

	t.Run("returns internal server error and logs error", func(t *testing.T) {
		t.Helper()

		request := statistico.TeamRequest{TeamId: 1}

		repo := new(mock.TeamRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)

		repo.On("ByID", uint64(1)).Return(&app.Team{}, errors.New("connection error"))

		_, err := service.GetTeamByID(context.Background(), &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = internal server error", err.Error())
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})
}

func TestTeamService_GetTeamsBySeasonId(t *testing.T) {
	t.Run("returns a slice of proto team struct", func(t *testing.T) {
		t.Helper()

		repo := new(mock.TeamRepository)
		logger, _ := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)
		server := new(mock.TeamServer)

		teams := []app.Team{
			newTeam(1, "West Ham United"),
			newTeam(2, "Arsenal"),
			newTeam(3, "Chelsea"),
		}

		repo.On("BySeasonId", uint64(16036)).Return(teams, nil)

		request := statistico.SeasonTeamsRequest{SeasonId: 16036}

		server.On("Send", mock2.AnythingOfType("*statistico.Team")).
			Times(3).
			Return(nil)

		err := service.GetTeamsBySeasonId(&request, server)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Nil(t, err)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned from team repository", func(t *testing.T) {
		t.Helper()

		repo := new(mock.TeamRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)
		server := new(mock.TeamServer)

		repo.On("BySeasonId", uint64(16036)).Return([]app.Team{}, errors.New("oh no"))

		server.AssertNotCalled(t, "Send")

		request := statistico.SeasonTeamsRequest{SeasonId: 16036}

		err := service.GetTeamsBySeasonId(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error retrieving Team(s) in Team Service. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned when streaming", func(t *testing.T) {
		t.Helper()

		repo := new(mock.TeamRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)
		server := new(mock.TeamServer)

		teams := []app.Team{
			newTeam(1, "West Ham United"),
			newTeam(2, "Arsenal"),
			newTeam(3, "Chelsea"),
		}

		repo.On("BySeasonId", uint64(16036)).Return(teams, nil)

		request := statistico.SeasonTeamsRequest{SeasonId: 16036}

		server.On("Send", mock2.AnythingOfType("*statistico.Team")).Return(errors.New("oh no"))

		err := service.GetTeamsBySeasonId(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error streaming Team back to client. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})
}

func newTeam(id uint64, name string) app.Team {
	return app.Team{
		ID:           id,
		Name:         name,
		VenueID:      560,
		CountryID:    uint64(462),
		NationalTeam: false,
		CreatedAt:    time.Unix(1546965200, 0),
		UpdatedAt:    time.Unix(1546965200, 0),
	}
}