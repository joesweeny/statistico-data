package grpc_test

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/errors"
	"github.com/statistico/statistico-data/internal/app/grpc"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTeamService_GetTeamById(t *testing.T) {
	t.Run("returns a proto Team struct", func(t *testing.T) {
		t.Helper()

		request := proto.TeamRequest{TeamId: 1}

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

		request := proto.TeamRequest{TeamId: 1}

		repo := new(mock.TeamRepository)
		logger, _ := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)

		repo.On("ByID", uint64(1)).Return(&app.Team{}, errors.ErrorNotFound)

		_, err := service.GetTeamByID(context.Background(), &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = NotFound desc = team with ID 1 does not exist", err.Error())
	})

	t.Run("returns internal server error and logs error", func(t *testing.T) {
		t.Helper()

		request := proto.TeamRequest{TeamId: 1}

		repo := new(mock.TeamRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewTeamService(repo, logger)

		repo.On("ByID", uint64(1)).Return(&app.Team{}, errors.ErrorDatabaseConnection)

		_, err := service.GetTeamByID(context.Background(), &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = internal server error", err.Error())
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})
}
