package grpc_test

import (
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTeamStatsService_GetStatForTeam(t *testing.T) {
	t.Run("returns a stream of proto team stat struct", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		statRepo := new(mock.TeamStatsRepository)
		xGRepo := new(mock.FixtureTeamXGRepository)
		logger, _ := test.NewNullLogger()

		service := grpc.NewTeamStatsService(fixtureRepo, statRepo, xGRepo, logger)

		server := new(mock.TeamStatServer)

		request := proto.TeamStatRequest{
			Stat:       "goals",
			TeamId:     55,
			SeasonIds:  []uint64{16036, 12968},
		}

		query := app.FixtureFilterQuery{SeasonIDs:  []uint64{16036, 12968}}

		fixtures := []app.Fixture{
			newFixture(1, 16036, 55, 2),
			newFixture(2, 16036, 49, 55),
		}

		fixtureRepo.On("ByTeamID", uint64(55), query).Return(fixtures, nil)

		goals := uint32(2)

		statOne := app.TeamStat{
			FixtureID: 1,
			Stat:      "goals",
			Value:     &goals,
		}

		statTwo := app.TeamStat{
			FixtureID: 2,
			Stat:      "goals",
			Value:     &goals,
		}

		statRepo.On("StatByFixtureAndTeam", "goals", uint64(1), uint64(55)).Return(&statOne, nil)
		statRepo.On("StatByFixtureAndTeam", "goals", uint64(2), uint64(55)).Return(&statTwo, nil)

		server.On("Send", mock2.AnythingOfType("*proto.TeamStat")).
			Times(2).
			Return(nil)

		err := service.GetStatForTeam(&request, server)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Nil(t, err)
		fixtureRepo.AssertExpectations(t)
		statRepo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("returns a stream of proto team stat struct handling opponent filter", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		statRepo := new(mock.TeamStatsRepository)
		xGRepo := new(mock.FixtureTeamXGRepository)
		logger, _ := test.NewNullLogger()

		service := grpc.NewTeamStatsService(fixtureRepo, statRepo, xGRepo, logger)

		server := new(mock.TeamStatServer)

		request := proto.TeamStatRequest{
			Stat:       "goals",
			TeamId:     55,
			SeasonIds:  []uint64{16036, 12968},
			Opponent:   &wrappers.BoolValue{Value: true},
		}

		query := app.FixtureFilterQuery{SeasonIDs:  []uint64{16036, 12968}}

		fixtures := []app.Fixture{
			newFixture(1, 16036, 55, 2),
			newFixture(2, 16036, 49, 55),
		}

		fixtureRepo.On("ByTeamID", uint64(55), query).Return(fixtures, nil)

		goals := uint32(2)

		statOne := app.TeamStat{
			FixtureID: 1,
			Stat:      "goals",
			Value:     &goals,
		}

		statTwo := app.TeamStat{
			FixtureID: 2,
			Stat:      "goals",
			Value:     &goals,
		}

		statRepo.On("StatByFixtureAndTeam", "goals", uint64(1), uint64(2)).Return(&statOne, nil)
		statRepo.On("StatByFixtureAndTeam", "goals", uint64(2), uint64(49)).Return(&statTwo, nil)

		server.On("Send", mock2.AnythingOfType("*proto.TeamStat")).
			Times(2).
			Return(nil)

		err := service.GetStatForTeam(&request, server)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Nil(t, err)
		fixtureRepo.AssertExpectations(t)
		statRepo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error fetching fixtures", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		statRepo := new(mock.TeamStatsRepository)
		xGRepo := new(mock.FixtureTeamXGRepository)
		logger, hook := test.NewNullLogger()

		service := grpc.NewTeamStatsService(fixtureRepo, statRepo, xGRepo, logger)

		server := new(mock.TeamStatServer)

		request := proto.TeamStatRequest{
			Stat:       "goals",
			TeamId:     55,
			SeasonIds:  []uint64{16036, 12968},
		}

		query := app.FixtureFilterQuery{SeasonIDs:  []uint64{16036, 12968}}

		fixtureRepo.On("ByTeamID", uint64(55), query).Return([]app.Fixture{}, errors.New("oh no"))

		statRepo.AssertNotCalled(t, "StatByFixtureAndTeam", "goals", uint64(1), uint64(2))

		server.AssertNotCalled(t, "Send", mock2.AnythingOfType("*proto.TeamStat"))

		err := service.GetStatForTeam(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error retrieving fixture(s) in team stats service. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
		fixtureRepo.AssertExpectations(t)
		statRepo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns empty stream if error fetching team stat", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		statRepo := new(mock.TeamStatsRepository)
		xGRepo := new(mock.FixtureTeamXGRepository)
		logger, hook := test.NewNullLogger()

		service := grpc.NewTeamStatsService(fixtureRepo, statRepo, xGRepo, logger)

		server := new(mock.TeamStatServer)

		request := proto.TeamStatRequest{
			Stat:       "goals",
			TeamId:     55,
			SeasonIds:  []uint64{16036, 12968},
		}

		query := app.FixtureFilterQuery{SeasonIDs:  []uint64{16036, 12968}}

		fixtures := []app.Fixture{
			newFixture(1, 16036, 55, 2),
		}

		fixtureRepo.On("ByTeamID", uint64(55), query).Return(fixtures, nil)

		statRepo.On("StatByFixtureAndTeam", "goals", uint64(1), uint64(55)).Return(&app.TeamStat{}, errors.New("oh no"))

		server.AssertNotCalled(t, "Send", mock2.AnythingOfType("*proto.TeamStat"))
		
		err := service.GetStatForTeam(&request, server)

		if err != nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "Error retrieving team stat in team stats service. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
		fixtureRepo.AssertExpectations(t)
		statRepo.AssertExpectations(t)
		server.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error stream team stats struct", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		statRepo := new(mock.TeamStatsRepository)
		xGRepo := new(mock.FixtureTeamXGRepository)
		logger, hook := test.NewNullLogger()

		service := grpc.NewTeamStatsService(fixtureRepo, statRepo, xGRepo, logger)

		server := new(mock.TeamStatServer)

		request := proto.TeamStatRequest{
			Stat:       "goals",
			TeamId:     55,
			SeasonIds:  []uint64{16036, 12968},
		}

		query := app.FixtureFilterQuery{SeasonIDs:  []uint64{16036, 12968}}

		fixtures := []app.Fixture{
			newFixture(1, 16036, 55, 2),
		}

		fixtureRepo.On("ByTeamID", uint64(55), query).Return(fixtures, nil)

		goals := uint32(2)

		statOne := app.TeamStat{
			FixtureID: 1,
			Stat:      "goals",
			Value:     &goals,
		}

		statRepo.On("StatByFixtureAndTeam", "goals", uint64(1), uint64(55)).Return(&statOne, nil)

		server.On("Send", mock2.AnythingOfType("*proto.TeamStat")).Return(errors.New("oh no"))

		err := service.GetStatForTeam(&request, server)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
		assert.Equal(t, "Error streaming team stat back to client. Error: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
		fixtureRepo.AssertExpectations(t)
		statRepo.AssertExpectations(t)
		server.AssertExpectations(t)
	})
}

func newFixture(id, seasonId, homeId, awayId uint64) app.Fixture {
	var roundId = uint64(165789)

	return app.Fixture{
		ID:         id,
		SeasonID:   seasonId,
		RoundID:    &roundId,
		HomeTeamID: homeId,
		AwayTeamID: awayId,
		Date:       time.Unix(1548086929, 0),
	}
}