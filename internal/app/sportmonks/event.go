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

func (e EventRequester) EventsByFixtureIDs(ids []uint64) (<-chan app.GoalEvent, <-chan app.SubstitutionEvent, <-chan app.CardEvent) {
	goal := make(chan app.GoalEvent, 500)
	sub := make(chan app.SubstitutionEvent, 500)
	card := make(chan app.CardEvent, 500)

	go e.parseEvents(ids, goal, sub, card)

	return goal, sub, card
}

func (e EventRequester) EventsBySeasonIDs(seasonIDs []uint64) (<-chan app.GoalEvent, <-chan app.SubstitutionEvent, <-chan app.CardEvent) {
	goal := make(chan app.GoalEvent, 500)
	sub := make(chan app.SubstitutionEvent, 500)
	card := make(chan app.CardEvent, 500)

	go e.parseBySeasonIDs(seasonIDs, goal, sub, card)

	return goal, sub, card
}

func (e EventRequester) parseBySeasonIDs(seasonIDs []uint64, g chan<- app.GoalEvent, s chan<- app.SubstitutionEvent, c chan<- app.CardEvent) {
	defer close(g)
	defer close(s)
	defer close(c)

	wg := sync.WaitGroup{}

	for _, id := range seasonIDs {
		wg.Add(1)
		go e.sendSeasonRequest(id, g, s, c, &wg)
	}

	wg.Wait()
}

func (e EventRequester) sendSeasonRequest(seasonID uint64, g chan<- app.GoalEvent, s chan<- app.SubstitutionEvent, c chan<- app.CardEvent, wg *sync.WaitGroup) {
	res, _, err := e.client.SeasonByID(context.Background(), int(seasonID), []string{"results.cards", "results.goals", "results.substitutions"})

	if err != nil {
		e.logger.Errorf(
			"Error when calling client '%s' when making season fixtures request. Season ID %d",
			err.Error(),
			seasonID,
		)

		wg.Done()
		return
	}

	for _, res := range res.Results() {
		for _, event := range res.Cards() {
			c <- transformCardEvent(event)
		}

		for _, event := range res.Goals() {
			g <- transformGoalEvent(event)
		}

		for _, event := range res.Substitutions() {
			s <- transformSubstitutionEvent(event)
		}
	}

	wg.Done()
}

func (e EventRequester) parseEvents(ids []uint64, g chan<- app.GoalEvent, s chan<- app.SubstitutionEvent, c chan<- app.CardEvent) {
	defer close(g)
	defer close(s)
	defer close(c)

	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)
		go e.sendEventRequest(id, g, s, c, &wg)
	}

	wg.Wait()
}

func (e EventRequester) sendEventRequest(fixtureId uint64, g chan<- app.GoalEvent, s chan<- app.SubstitutionEvent, c chan<- app.CardEvent, w *sync.WaitGroup) {
	includes := []string{"cards", "goals", "substitutions"}
	var filters map[string][]int

	res, _, err := e.client.FixtureByID(context.Background(), int(fixtureId), includes, filters)

	if err != nil {
		e.logger.Errorf(
			"Error when calling client '%s' when making fixture event request. Fixture ID %d",
			err.Error(),
			fixtureId,
		)
		w.Done()
		return
	}

	for _, event := range res.Cards() {
		c <- transformCardEvent(event)
	}

	for _, event := range res.Goals() {
		g <- transformGoalEvent(event)
	}

	for _, event := range res.Substitutions() {
		s <- transformSubstitutionEvent(event)
	}

	w.Done()
}

func transformCardEvent(s spClient.CardEvent) app.CardEvent {
	teamId, _ := strconv.Atoi(s.TeamID)

	return app.CardEvent{
		ID:          uint64(s.ID),
		TeamID:      uint64(teamId),
		Type:        s.Type,
		FixtureID:   uint64(s.FixtureID),
		PlayerID:    uint64(s.PlayerID),
		Minute:      uint8(s.Minute),
		Reason:      s.Reason,
	}
}

func transformGoalEvent(s spClient.GoalEvent) app.GoalEvent {
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

	return event
}

func transformSubstitutionEvent(s spClient.SubstitutionEvent) app.SubstitutionEvent {
	teamId, _ := strconv.Atoi(s.TeamID)

	return app.SubstitutionEvent{
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
