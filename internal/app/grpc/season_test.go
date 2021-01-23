package grpc_test

import (
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
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

func TestSeasonService_GetSeasonsForCompetition(t *testing.T) {
	t.Run("returns a stream of proto competition struct", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		logger, _ := test.NewNullLogger()
		service := grpc.NewSeasonService(repo, logger)
		server := new(mock.SeasonServer)

		seasons := []app.Season{
			newSeason(1, 16036, "2017/2018", false),
			newSeason(2, 16036, "2018/2019", false),
			newSeason(3, 16036, "2019/2020", true),
		}

		repo.On("ByCompetitionId", uint64(16036), "name_asc").Return(seasons, nil)

		request := statistico.SeasonCompetitionRequest{
			CompetitionId: 16036,
			Sort:          &wrappers.StringValue{Value: "name_asc"},
		}

		server.On("Send", mock2.AnythingOfType("*statistico.Season")).
			Times(3).
			Return(nil)

		err := service.GetSeasonsForCompetition(&request, server)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Nil(t, err)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned from season repository", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewSeasonService(repo, logger)
		server := new(mock.SeasonServer)

		repo.On("ByCompetitionId", uint64(16036), "name_asc").Return([]app.Season{}, errors.New("oh no"))

		server.AssertNotCalled(t, "Send")

		request := statistico.SeasonCompetitionRequest{
			CompetitionId: 16036,
			Sort:          &wrappers.StringValue{Value: "name_asc"},
		}

		err := service.GetSeasonsForCompetition(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error retrieving Season(s) in Season Service. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned when streaming response", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewSeasonService(repo, logger)
		server := new(mock.SeasonServer)

		seasons := []app.Season{
			newSeason(1, 16036, "2017/2018", false),
			newSeason(2, 16036, "2018/2019", false),
			newSeason(3, 16036, "2019/2020", true),
		}

		repo.On("ByCompetitionId", uint64(16036), "name_asc").Return(seasons, nil)

		request := statistico.SeasonCompetitionRequest{
			CompetitionId: 16036,
			Sort:          &wrappers.StringValue{Value: "name_asc"},
		}

		server.On("Send", mock2.AnythingOfType("*statistico.Season")).Return(errors.New("oh no"))

		err := service.GetSeasonsForCompetition(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error streaming Season back to client. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})
}

func TestSeasonService_GetSeasonsForTeam(t *testing.T) {
	t.Run("returns a stream of proto competition struct", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		logger, _ := test.NewNullLogger()
		service := grpc.NewSeasonService(repo, logger)
		server := new(mock.SeasonServer)

		seasons := []app.Season{
			newSeason(1, 16036, "2017/2018", false),
			newSeason(2, 16036, "2018/2019", false),
			newSeason(3, 16036, "2019/2020", true),
		}

		repo.On("ByTeamId", uint64(1), "name_asc").Return(seasons, nil)

		request := statistico.TeamSeasonsRequest{
			TeamId: 	   1,
			Sort:          &wrappers.StringValue{Value: "name_asc"},
		}

		server.On("Send", mock2.AnythingOfType("*statistico.Season")).
			Times(3).
			Return(nil)

		err := service.GetSeasonsForTeam(&request, server)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Nil(t, err)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned from season repository", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewSeasonService(repo, logger)
		server := new(mock.SeasonServer)

		repo.On("ByTeamId", uint64(1), "name_asc").Return([]app.Season{}, errors.New("oh no"))

		server.AssertNotCalled(t, "Send")

		request := statistico.TeamSeasonsRequest{
			TeamId: 	   1,
			Sort:          &wrappers.StringValue{Value: "name_asc"},
		}

		err := service.GetSeasonsForTeam(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error retrieving Season(s) in Season Service. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error returned when streaming response", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		logger, hook := test.NewNullLogger()
		service := grpc.NewSeasonService(repo, logger)
		server := new(mock.SeasonServer)

		seasons := []app.Season{
			newSeason(1, 16036, "2017/2018", false),
			newSeason(2, 16036, "2018/2019", false),
			newSeason(3, 16036, "2019/2020", true),
		}

		repo.On("ByTeamId", uint64(16036), "name_asc").Return(seasons, nil)

		request := statistico.TeamSeasonsRequest{
			TeamId: 	   16036,
			Sort:          &wrappers.StringValue{Value: "name_asc"},
		}

		server.On("Send", mock2.AnythingOfType("*statistico.Season")).Return(errors.New("oh no"))

		err := service.GetSeasonsForTeam(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error streaming Season back to client. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		repo.AssertExpectations(t)
		server.AssertExpectations(t)
	})
}

func newSeason(id uint64, competitionId uint64, name string, current bool) app.Season {
	return app.Season{
		ID:            id,
		Name:          name,
		CompetitionID: competitionId,
		IsCurrent:     current,
		CreatedAt:     time.Unix(1546965200, 0),
		UpdatedAt:     time.Unix(1546965200, 0),
	}
}
