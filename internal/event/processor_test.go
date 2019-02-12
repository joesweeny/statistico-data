package event

import (
	"github.com/stretchr/testify/mock"
	"github.com/joesweeny/statshub/internal/model"
	"testing"
	"github.com/jonboulle/clockwork"
	"log"
	"io/ioutil"
)

func TestProcessGoalEvent(t *testing.T) {
	repo := new(mockRepository)
	processor := Processor{
		Repository: repo,
		Factory:    Factory{clockwork.NewFakeClock()},
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	event := newClientGoalEvent()

	t.Run("creates and inserts new goal event if not present in database", func(t *testing.T) {
		repo.On("GoalEventById", 55).Return(&model.GoalEvent{}, ErrNotFound)
		repo.On("InsertGoalEvent", mock.Anything).Return(nil)
		processor.ProcessGoalEvent(event)
	})

	t.Run("goal event is not inserted if it already exists", func(t *testing.T) {
		repo.On("GoalEventById", 55).Return(&model.GoalEvent{}, nil)
		repo.AssertNotCalled(t,"InsertGoalEvent", mock.Anything)
		processor.ProcessGoalEvent(event)
	})
}

func TestProcessSubstitutionEvent(t *testing.T) {
	repo := new(mockRepository)
	processor := Processor{
		Repository: repo,
		Factory:    Factory{clockwork.NewFakeClock()},
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	event := newClientSubstitutionEvent()

	t.Run("creates and inserts new sub event if not present in database", func(t *testing.T) {
		repo.On("SubstitutionEventById", 57).Return(&model.SubstitutionEvent{}, ErrNotFound)
		repo.On("InsertSubstitutionEvent", mock.Anything).Return(nil)
		processor.ProcessSubstitutionEvent(event)
	})

	t.Run("sub event is not inserted if it already exists", func(t *testing.T) {
		repo.On("SubstitutionEventById", 57).Return(&model.SubstitutionEvent{}, nil)
		repo.AssertNotCalled(t,"InsertSubstitutionEvent", mock.Anything)
		processor.ProcessSubstitutionEvent(event)
	})
}

type mockRepository struct {
	mock.Mock
}

func (m mockRepository) InsertGoalEvent(e *model.GoalEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m mockRepository) InsertSubstitutionEvent(e *model.SubstitutionEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m mockRepository) GoalEventById(id int) (*model.GoalEvent, error) {
	args := m.Called(id)
	return args.Get(0).(*model.GoalEvent), args.Error(1)
}

func (m mockRepository) SubstitutionEventById(id int) (*model.SubstitutionEvent, error) {
	args := m.Called(id)
	return args.Get(0).(*model.SubstitutionEvent), args.Error(1)
}


