package grpc_test

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventService_FixtureEvents(t *testing.T) {
	t.Run("returns a proto FixtureEventsResponse struct", func(t *testing.T) {
		t.Helper()

		repo := new(mock.EventRepository)
		logger, _ := test.NewNullLogger()

		service := grpc.NewEventService(repo, logger)

		cards := []*app.CardEvent{
			newCardEvent(1),
			newCardEvent(2),
		}

		goals := []*app.GoalEvent{
			newGoalEvent(3),
			newGoalEvent(4),
		}

		repo.On("CardEventsForFixture", uint64(45)).Return(cards, nil)
		repo.On("GoalEventsForFixture", uint64(45)).Return(goals, nil)

		req := statistico.FixtureRequest{FixtureId: uint64(45)}

		res, err := service.FixtureEvents(context.Background(), &req)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err)
		}

		expectedCards := []*statistico.CardEvent{
			{
				Id: 1,
				TeamId: 4509,
				Type: "redcard",
				PlayerId: 3401,
				Minute: 85,
			},
			{
				Id: 2,
				TeamId: 4509,
				Type: "redcard",
				PlayerId: 3401,
				Minute: 85,
			},
		}

		expectedGoals := []*statistico.GoalEvent{
			{
				Id: 3,
				TeamId: 4509,
				PlayerId: 3401,
				PlayerAssistId: nil,
				Minute: 82,
				Score: "0-1",
			},
			{
				Id: 4,
				TeamId: 4509,
				PlayerId: 3401,
				PlayerAssistId: nil,
				Minute: 82,
				Score: "0-1",
			},
		}

		assert.Equal(t, uint64(45), res.FixtureId)
		assert.Equal(t, expectedCards, res.Cards)
		assert.Equal(t, expectedGoals, res.Goals)
		repo.AssertExpectations(t)
	})

	t.Run("logs error if error returned by repository when fetching card events", func(t *testing.T) {
		t.Helper()

		repo := new(mock.EventRepository)
		logger, hook := test.NewNullLogger()

		service := grpc.NewEventService(repo, logger)

		repo.On("CardEventsForFixture", uint64(45)).Return([]*app.CardEvent{}, errors.New("oh no"))
		repo.AssertNotCalled(t, "GoalEventsForFixture", uint64(45))

		req := statistico.FixtureRequest{FixtureId: uint64(45)}

		_, err := service.FixtureEvents(context.Background(), &req)

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Equal(t, "Error retrieving card events in Event Service. Error: oh no", hook.LastEntry().Message)
		repo.AssertExpectations(t)
	})

	t.Run("logs error if error returned by repository when fetching goal events", func(t *testing.T) {
		t.Helper()

		repo := new(mock.EventRepository)
		logger, hook := test.NewNullLogger()

		service := grpc.NewEventService(repo, logger)

		repo.On("CardEventsForFixture", uint64(45)).Return([]*app.CardEvent{}, nil)
		repo.On("GoalEventsForFixture", uint64(45)).Return([]*app.GoalEvent{}, errors.New("oh no"))

		req := statistico.FixtureRequest{FixtureId: uint64(45)}

		_, err := service.FixtureEvents(context.Background(), &req)

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Equal(t, "Error retrieving goal events in Event Service. Error: oh no", hook.LastEntry().Message)
		repo.AssertExpectations(t)
	})
}

func newCardEvent(id uint64) *app.CardEvent {
	return &app.CardEvent{
		ID:          id,
		TeamID:      uint64(4509),
		Type:        "redcard",
		FixtureID:   uint64(45),
		PlayerID:    uint64(3401),
		Minute:      85,
		Reason:      nil,
		CreatedAt:   time.Unix(1546965200, 0),
	}
}

func newGoalEvent(id uint64) *app.GoalEvent {
	return &app.GoalEvent{
		ID:        id,
		FixtureID: uint64(45),
		TeamID:    uint64(4509),
		PlayerID:  uint64(3401),
		Minute:    82,
		Score:     "0-1",
		CreatedAt: time.Unix(1546965200, 0),
	}
}
