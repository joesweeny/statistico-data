package grpc_test

import (
	"github.com/golang/protobuf/ptypes/wrappers"
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

		request := proto.SeasonCompetitionRequest{
			CompetitionId: 16036,
			Sort:          &wrappers.StringValue{Value: "name_asc"},
		}

		server.On("Send", mock2.AnythingOfType("*proto.Season")).
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
