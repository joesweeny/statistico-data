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
	t.Run("inserts new events into repository when processing events by result id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		fixtureRepo.On("ByID", uint64(12)).Return(newFixture(uint64(12)), nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", goalOne).Return(nil)
		eventRepo.On("InsertGoalEvent", goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", subTwo).Return(nil)

		processor.Process("events:by-result-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("does not insert events into repository when processing events by result id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		fixtureRepo.On("ByID", uint64(12)).Return(newFixture(uint64(12)), nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(goalOne, nil)
		eventRepo.On("GoalEventByID", uint64(20)).Return(goalTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", goalOne)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", goalTwo)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(subOne, nil)
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(subTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", subOne)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", subTwo)

		processor.Process("events:by-result-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert event into repository when processing events by result id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		fixtureRepo.On("ByID", uint64(12)).Return(newFixture(uint64(12)), nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", goalOne).Return(errors.New("error occurred"))
		eventRepo.On("InsertGoalEvent", goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", subTwo).Return(errors.New("error occurred"))

		processor.Process("events:by-result-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Equal(t, 2, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new events into repository when processing events by season id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		var fix []app.Fixture

		fix = append(fix, *newFixture(uint64(12)))

		query := app.FixtureRepositoryQuery{SeasonIDs: []uint64{uint64(12)}}

		fixtureRepo.On("Get", query).Return(fix, nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", goalOne).Return(nil)
		eventRepo.On("InsertGoalEvent", goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", subTwo).Return(nil)

		processor.Process("events:by-season-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("does not insert events into repository when processing events by season id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		var fix []app.Fixture
		fix = append(fix, *newFixture(uint64(12)))

		query := app.FixtureRepositoryQuery{SeasonIDs: []uint64{uint64(12)}}

		fixtureRepo.On("Get", query).Return(fix, nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(goalOne, nil)
		eventRepo.On("GoalEventByID", uint64(20)).Return(goalTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", goalOne)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", goalTwo)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(subOne, nil)
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(subTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", subOne)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", subTwo)

		processor.Process("events:by-season-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert event into repository when processing events by season id command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		var fix []app.Fixture
		fix = append(fix, *newFixture(uint64(12)))

		query := app.FixtureRepositoryQuery{SeasonIDs: []uint64{uint64(12)}}

		fixtureRepo.On("Get", query).Return(fix, nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", goalOne).Return(errors.New("error occurred"))
		eventRepo.On("InsertGoalEvent", goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", subTwo).Return(errors.New("error occurred"))

		processor.Process("events:by-season-id", "12", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Equal(t, 2, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new events into repository when processing events today command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		query := app.FixtureRepositoryQuery{DateFrom: &from, DateTo: &to}

		fixtureRepo.On("GetIDs", query).Return([]uint64{12}, nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", goalOne).Return(nil)
		eventRepo.On("InsertGoalEvent", goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", subTwo).Return(nil)

		processor.Process("events:today", "", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("does not insert events into repository when processing events today command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		query := app.FixtureRepositoryQuery{DateFrom: &from, DateTo: &to}

		fixtureRepo.On("GetIDs", query).Return([]uint64{12}, nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(goalOne, nil)
		eventRepo.On("GoalEventByID", uint64(20)).Return(goalTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", goalOne)
		eventRepo.AssertNotCalled(t, "InsertGoalEvent", goalTwo)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(subOne, nil)
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(subTwo, nil)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", subOne)
		eventRepo.AssertNotCalled(t, "InsertSubstitutionEvent", subTwo)

		processor.Process("events:today", "", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert event into repository when processing events today command", func(t *testing.T) {
		t.Helper()

		eventRepo := new(mock.EventRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.EventRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewEventProcessor(eventRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		goalOne := newGoalEvent(10)
		goalTwo := newGoalEvent(20)

		goals := make([]*app.GoalEvent, 2)
		goals[0] = goalOne
		goals[1] = goalTwo

		goalCh := goalEventChannel(goals)

		subOne := newSubstitutionEvent(55)
		subTwo := newSubstitutionEvent(2)

		subs := make([]*app.SubstitutionEvent, 2)
		subs[0] = subOne
		subs[1] = subTwo

		subCh := subEventChannel(subs)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		query := app.FixtureRepositoryQuery{DateFrom: &from, DateTo: &to}

		fixtureRepo.On("GetIDs", query).Return([]uint64{12}, nil)
		requester.On("EventsByFixtureIDs", []uint64{12}).Return(goalCh, subCh)

		eventRepo.On("GoalEventByID", uint64(10)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("GoalEventByID", uint64(20)).Return(&app.GoalEvent{}, errors.New("not found"))
		eventRepo.On("InsertGoalEvent", goalOne).Return(errors.New("error occurred"))
		eventRepo.On("InsertGoalEvent", goalTwo).Return(nil)

		eventRepo.On("SubstitutionEventByID", uint64(55)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("SubstitutionEventByID", uint64(2)).Return(&app.SubstitutionEvent{}, errors.New("not found"))
		eventRepo.On("InsertSubstitutionEvent", subOne).Return(nil)
		eventRepo.On("InsertSubstitutionEvent", subTwo).Return(errors.New("error occurred"))

		processor.Process("events:today", "", done)

		<-done

		requester.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		eventRepo.AssertExpectations(t)
		assert.Equal(t, 2, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
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

func newSubstitutionEvent(id uint64) *app.SubstitutionEvent {
	true := true
	return &app.SubstitutionEvent{
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

func goalEventChannel(goals []*app.GoalEvent) chan *app.GoalEvent {
	ch := make(chan *app.GoalEvent, len(goals))

	for _, c := range goals {
		ch <- c
	}

	close(ch)

	return ch
}

func subEventChannel(subs []*app.SubstitutionEvent) chan *app.SubstitutionEvent {
	ch := make(chan *app.SubstitutionEvent, len(subs))

	for _, c := range subs {
		ch <- c
	}

	close(ch)

	return ch
}
