package process_test

import (
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/process"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventProcessor_Process(t *testing.T) {
	t.Run("inserts new events into repository when processing events by season id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		cardOne := newCardEvent(1)
		cardTwo := newCardEvent(2)

		cards := make([]app.CardEvent, 2)
		cards[0] = cardOne
		cards[1] = cardTwo

		cardCh := cardEventChannel(cards)

		requester.On("EventsBySeasonIDs", []uint64{12}).Return(goalCh, subCh, cardCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", &goalOne).Return(nil)
		eventRepo.On("InsertGoalEvent", &goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", &subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", &subTwo).Return(nil)

		eventRepo.On("CardEventByID", uint64(1)).Return(&app.CardEvent{}, errors.New("not found"))
		eventRepo.On("CardEventByID", uint64(2)).Return(&app.CardEvent{}, errors.New("not found"))
		eventRepo.On("InsertCardEvent", &cardOne).Return(nil)
		eventRepo.On("InsertCardEvent", &cardTwo).Return(nil)

		processor.Process("events:by-season-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("does not insert events into repository when processing events by season id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		cardOne := newCardEvent(1)
		cardTwo := newCardEvent(2)

		cards := make([]app.CardEvent, 2)
		cards[0] = cardOne
		cards[1] = cardTwo

		cardCh := cardEventChannel(cards)

		requester.On("EventsBySeasonIDs", []uint64{12}).Return(goalCh, subCh, cardCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&goalOne, nil)
		eventRepo.On("GoalEventByID", uint64(20)).Return(&goalTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", &goalOne)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", &goalTwo)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&subOne, nil)
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&subTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", &subOne)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", &subTwo)

		eventRepo.On("CardEventByID", uint64(1)).Return(&cardOne, nil)
		eventRepo.On("CardEventByID", uint64(2)).Return(&cardTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertCardEvent", &cardOne)
		eventRepo.AssertNotCalled(t, "InsertCardEvent", &cardTwo)

		processor.Process("events:by-season-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert event into repository when processing events by season id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		cardOne := newCardEvent(1)
		cardTwo := newCardEvent(2)

		cards := make([]app.CardEvent, 2)
		cards[0] = cardOne
		cards[1] = cardTwo

		cardCh := cardEventChannel(cards)

		requester.On("EventsBySeasonIDs", []uint64{12}).Return(goalCh, subCh, cardCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", &goalOne).Return(errors.New("error occurred"))
		eventRepo.On("InsertGoalEvent", &goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", &subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", &subTwo).Return(errors.New("error occurred"))

		eventRepo.On("CardEventByID", uint64(1)).Return(&app.CardEvent{}, errors.New("not found"))
		eventRepo.On("CardEventByID", uint64(2)).Return(&app.CardEvent{}, errors.New("not found"))
		eventRepo.On("InsertCardEvent", &cardOne).Return(nil)
		eventRepo.On("InsertCardEvent", &cardTwo).Return(errors.New("error occurred"))

		processor.Process("events:by-season-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Equal(t, 3, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newCardEvent(id uint64) app.CardEvent {
	return app.CardEvent{
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

func newGoalEvent(id uint64) app.GoalEvent {
	return app.GoalEvent{
		ID:        id,
		FixtureID: uint64(45),
		TeamID:    uint64(4509),
		PlayerID:  uint64(3401),
		Minute:    82,
		Score:     "0-1",
		CreatedAt: time.Unix(1546965200, 0),
	}
}

func newSubstitutionEvent(id uint64) app.SubstitutionEvent {
	true := true
	return app.SubstitutionEvent{
		ID:          id,
		FixtureID:   uint64(45),
		TeamID:      uint64(4509),
		PlayerInID:  uint64(3401),
		PlayerOutID: uint64(901),
		Minute:      82,
		Injured:     &true,
		CreatedAt:   time.Unix(1546965200, 0),
	}
}

func cardEventChannel(cards []app.CardEvent) chan app.CardEvent {
	ch := make(chan app.CardEvent, len(cards))

	for _, c := range cards {
		ch <- c
	}

	close(ch)

	return ch
}

func goalEventChannel(goals []app.GoalEvent) chan app.GoalEvent {
	ch := make(chan app.GoalEvent, len(goals))

	for _, c := range goals {
		ch <- c
	}

	close(ch)

	return ch
}

func subEventChannel(subs []app.SubstitutionEvent) chan app.SubstitutionEvent {
	ch := make(chan app.SubstitutionEvent, len(subs))

	for _, c := range subs {
		ch <- c
	}

	close(ch)

	return ch
}
