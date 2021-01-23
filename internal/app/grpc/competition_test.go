package grpc_test

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCompetitionService_ListCompetitions(t *testing.T) {
	t.Run("returns a stream of proto competition struct", func(t *testing.T) {
		t.Helper()

		repo := new(mock.CompetitionRepository)
		logger, _ := test.NewNullLogger()
		service := grpc.NewCompetitionService(repo, logger)
		server := new(mock.CompetitionServer)

		competitions := []app.Competition{
			newCompetition(8, 462, false),
			newCompetition(12, 462, false),
			newCompetition(16, 462, false),
		}

		query := app.CompetitionFilterQuery{
			CountryIds: []uint64{462},
		}

		repo.On("Get", query).Return(competitions, nil)

		request := statistico.CompetitionRequest{
			CountryIds: []uint64{462},
		}

		server.On("Send", mock2.AnythingOfType("*statistico.Competition")).
			Times(3).
			Return(nil)

		err := service.ListCompetitions(&request, server)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Nil(t, err)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned from competition repository", func(t *testing.T) {
		t.Helper()

		repo := new(mock.CompetitionRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewCompetitionService(repo, logger)
		server := new(mock.CompetitionServer)

		repo.On("Get", app.CompetitionFilterQuery{}).Return([]app.Competition{}, errors.New("oh no"))

		server.AssertNotCalled(t, "Send")

		err := service.ListCompetitions(&statistico.CompetitionRequest{}, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error retrieving Competition(s) in Competition Service. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned when streaming response", func(t *testing.T) {
		t.Helper()

		repo := new(mock.CompetitionRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewCompetitionService(repo, logger)
		server := new(mock.CompetitionServer)

		competitions := []app.Competition{
			newCompetition(8, 462, false),
			newCompetition(12, 462, false),
			newCompetition(16, 462, false),
		}

		query := app.CompetitionFilterQuery{
			CountryIds: []uint64{462},
		}

		repo.On("Get", query).Return(competitions, nil)

		request := statistico.CompetitionRequest{
			CountryIds: []uint64{462},
		}

		server.On("Send", mock2.AnythingOfType("*statistico.Competition")).Return(errors.New("oh no"))

		err := service.ListCompetitions(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error streaming Competition back to client. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})
}

func newCompetition(id uint64, country uint64, isCup bool) app.Competition {
	return app.Competition{
		ID:        id,
		Name:      "Premier League",
		CountryID: country,
		IsCup:     isCup,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}
