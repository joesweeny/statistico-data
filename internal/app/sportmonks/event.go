package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/helpers"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"strconv"
	"sync"
)

type EventRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (e EventRequester) EventsByFixtureIDs(ids []uint64) (<-chan *app.GoalEvent, <-chan *app.SubstitutionEvent) {
	goal := make(chan *app.GoalEvent, 500)
	sub := make(chan *app.SubstitutionEvent, 500)

	go e.parseEvents(ids, goal, sub)

	return goal, sub
}

func (e EventRequester) parseEvents(ids []uint64, g chan<- *app.GoalEvent, s chan<- *app.SubstitutionEvent) {
	defer close(g)
	defer close(s)

	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)
		go e.sendEventRequest(id, g, s, &wg)
	}

	wg.Wait()
}

func (e EventRequester) sendEventRequest(fixtureId uint64, g chan<- *app.GoalEvent, s chan<- *app.SubstitutionEvent, w *sync.WaitGroup) {
	includes := []string{"goals", "substitutions"}
	var filters map[string][]int

	res, _, err := e.client.FixtureByID(context.Background(), int(fixtureId), includes, filters)

	if err != nil {
		e.logger.Warnf(
			"Error when calling client '%s' when making fixture event request. Fixture ID %d",
			err.Error(),
			fixtureId,
		)
		w.Done()
		return
	}

	for _, event := range res.Goals() {
		g <- transformGoalEvent(&event)
	}

	for _, event := range res.Substitutions() {
		s <- transformSubstitutionEvent(&event)
	}

	w.Done()
}

func transformGoalEvent(s *spClient.GoalEvent) *app.GoalEvent {
	teamId, _ := strconv.Atoi(s.TeamID)

	event := app.GoalEvent{
		ID:             uint64(s.ID),
		FixtureID:      uint64(s.FixtureID),
		TeamID:         uint64(teamId),
		PlayerID:       uint64(s.PlayerID),
		PlayerAssistID: helpers.NullableUint64(s.PlayerAssistID),
		Minute:         s.Minute,
		Score:          s.Result,
	}

	return &event
}

func transformSubstitutionEvent(s *spClient.SubstitutionEvent) *app.SubstitutionEvent {
	teamId, _ := strconv.Atoi(s.TeamID)

	return &app.SubstitutionEvent{
		ID:          uint64(s.ID),
		FixtureID:   uint64(s.FixtureID),
		TeamID:      uint64(teamId),
		PlayerInID:  uint64(s.PlayerInID),
		PlayerOutID: uint64(s.PlayerOutID),
		Minute:      s.Minute,
		Injured:     s.Injured,
	}
}

func NewEventRequester(client *spClient.HTTPClient, log *logrus.Logger) *EventRequester {
	return &EventRequester{client: client, logger: log}
}
