package grpc_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-data/internal/app/grpc"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/performance"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPerformanceService_GetTeamsMatchingStat(t *testing.T) {
	request := proto.TeamStatPerformanceRequest{
		Action:  "for",
		Games:   3,
		Measure: "average",
		Metric:  "gte",
		Seasons: []uint64{16036},
		Stat:    "goals",
		Value:   2.5,
		Venue:   "home",
	}

	filter := performance.StatFilter{
		Action:  "for",
		Games:   3,
		Measure: "average",
		Metric:  "gte",
		Seasons: []uint64{16036},
		Stat:    "goals",
		Value:   2.5,
		Venue:   "home",
	}

	t.Run("returns a TeamStatResponse struct containing team information", func(t *testing.T) {
		t.Helper()

		teams := []*performance.Team{
			{
				ID:   1,
				Name: "West Ham United",
			},
			{
				ID:   8,
				Name: "Liverpool",
			},
		}

		reader := new(mock.StatReader)
		service := grpc.NewPerformanceService(reader)

		reader.On("TeamsMatchingFilter", &filter).Return(teams, nil)

		response, err := service.GetTeamsMatchingStat(context.Background(), &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		protoTeams := []*proto.Team{
			{
				Id:   1,
				Name: "West Ham United",
			},
			{
				Id:   8,
				Name: "Liverpool",
			},
		}

		assert.Equal(t, 2, len(response.Teams))
		assert.Equal(t, protoTeams, response.Teams)
	})

	t.Run("returns error if error returned by reader", func(t *testing.T) {
		t.Helper()

		reader := new(mock.StatReader)
		service := grpc.NewPerformanceService(reader)

		reader.On("TeamsMatchingFilter", &filter).Return([]*performance.Team{}, errors.New("error occurred"))

		_, err := service.GetTeamsMatchingStat(context.Background(), &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rpc error: code = Internal desc = Internal server error", err.Error())
	})
}