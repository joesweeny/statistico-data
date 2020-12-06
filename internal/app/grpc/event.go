package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-proto/data/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventService struct {
	eventRepo  app.EventRepository
	logger     *logrus.Logger
}

func (e *EventService) FixtureEvents(ctx context.Context, req *statisticoproto.FixtureRequest) (*statisticoproto.FixtureEventsResponse, error) {
	cards, err := e.eventRepo.CardEventsForFixture(req.FixtureId)

	if err != nil {
		e.logger.Errorf("Error retrieving card events in Event Service. Error: %s", err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	goals, err := e.eventRepo.GoalEventsForFixture(req.FixtureId)

	if err != nil {
		e.logger.Errorf("Error retrieving goal events in Event Service. Error: %s", err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	res := statisticoproto.FixtureEventsResponse{FixtureId: req.FixtureId}

	for _, c := range cards {
		res.Cards = append(res.Cards, factory.CardEventToProto(c))
	}

	for _, g := range goals {
		res.Goals = append(res.Goals, factory.GoalEventToProto(g))
	}

	return &res, nil
}

func NewEventService(r app.EventRepository, l *logrus.Logger) *EventService {
	return &EventService{
		eventRepo: r,
		logger:    l,
	}
}
