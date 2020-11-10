package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type EventRepository struct {
	mock.Mock
}

func (m *EventRepository) InsertCardEvent(e *app.CardEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *EventRepository) InsertGoalEvent(e *app.GoalEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *EventRepository) InsertSubstitutionEvent(e *app.SubstitutionEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *EventRepository) CardEventsForFixture(fixtureID uint64) ([]*app.CardEvent, error) {
	args := m.Called(fixtureID)
	return args.Get(0).([]*app.CardEvent), args.Error(1)
}

func (m *EventRepository) GoalEventsForFixture(fixtureID uint64) ([]*app.GoalEvent, error) {
	args := m.Called(fixtureID)
	return args.Get(0).([]*app.GoalEvent), args.Error(1)
}

func (m *EventRepository) GoalEventByID(id uint64) (*app.GoalEvent, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.GoalEvent)
	return c, args.Error(1)
}

func (m *EventRepository) SubstitutionEventByID(id uint64) (*app.SubstitutionEvent, error) {
	args := m.Called(id)
	c := args.Get(0).(*app.SubstitutionEvent)
	return c, args.Error(1)
}

type EventRequester struct {
	mock.Mock
}

func (m *EventRequester) EventsByFixtureIDs(ids []uint64) (<-chan *app.GoalEvent, <-chan *app.SubstitutionEvent, <-chan *app.CardEvent) {
	args := m.Called(ids)
	return args.Get(0).(chan *app.GoalEvent), args.Get(1).(chan *app.SubstitutionEvent), args.Get(2).(chan *app.CardEvent)
}
