package sportmonks

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"strconv"
	"sync"
)

type EventRequester struct {
	client *spClient.HTTPClient
	logger *logrus.Logger
}

func (e EventRequester) EventsByFixtureIDs(ids []int64) (<-chan *app.GoalEvent, <-chan *app.SubstitutionEvent) {
	goal := make(chan *app.GoalEvent, 500)
	sub := make(chan *app.SubstitutionEvent, 500)

	go e.parseEvents(ids, goal, sub)

	return goal, sub
}

func (e EventRequester) parseEvents(ids []int64, g chan<- *app.GoalEvent, s chan<- *app.SubstitutionEvent) {
	defer close(g)
	defer close(s)

	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)
		go e.sendEventRequest(id, g, s, &wg)
	}

	wg.Wait()
}

func (e EventRequester) sendEventRequest(fixtureId int64, g chan<- *app.GoalEvent, s chan<- *app.SubstitutionEvent, w *sync.WaitGroup) {
	includes := []string{"goals", "substitutions"}
	var filters map[string][]int

	res, _, err := e.client.FixtureByID(context.Background(), int(fixtureId), includes, filters)

	if err != nil {
		e.logger.Fatalf("Error when calling client '%s' when making fixture event request", err.Error())
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
	assist := int64(*s.PlayerAssistID)

	return &app.GoalEvent{
		ID:             int64(s.ID),
		FixtureID:      int64(s.FixtureID),
		TeamID:         int64(teamId),
		PlayerID:       int64(s.PlayerID),
		PlayerAssistID: &assist,
		Minute:         s.Minute,
		Score:          s.Result,
	}
}

func transformSubstitutionEvent(s *spClient.SubstitutionEvent) *app.SubstitutionEvent {
	teamId, _ := strconv.Atoi(s.TeamID)

	return &app.SubstitutionEvent{
		ID:             int64(s.ID),
		FixtureID:      int64(s.FixtureID),
		TeamID:         int64(teamId),
		PlayerInID:  int64(s.PlayerInID),
		PlayerOutID: int64(s.PlayerOutID),
		Minute:      s.Minute,
		Injured:     s.Injured,
	}
}
