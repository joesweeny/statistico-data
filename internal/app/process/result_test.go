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

func TestResultProcessor_Process(t *testing.T) {
	t.Run("inserts new result into repository when processing result by id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fixtureRepo.On("ByID", uint64(34)).Return(newFixture(34), nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(&app.Result{}, errors.New("not found"))
		resultRepo.On("Insert", res).Return(nil)
		processor.Process("result:by-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing result into repository when processing result by id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fixtureRepo.On("ByID", uint64(34)).Return(newFixture(34), nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(res, nil)
		resultRepo.On("Update", res).Return(nil)
		processor.Process("result:by-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert result into repository when processing result by id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fixtureRepo.On("ByID", uint64(34)).Return(newFixture(34), nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(&app.Result{}, errors.New("not found"))
		resultRepo.On("Insert", res).Return(errors.New("error occurred"))
		processor.Process("result:by-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update result into repository when processing result by id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fixtureRepo.On("ByID", uint64(34)).Return(newFixture(34), nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(res, nil)
		resultRepo.On("Update", res).Return(errors.New("error occurred"))
		processor.Process("result:by-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new result into repository when processing result by season id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fix := []app.Fixture{*newFixture(34)}

		fixtureRepo.On("BySeasonID", uint64(34)).Return(fix, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(&app.Result{}, errors.New("not found"))
		resultRepo.On("Insert", res).Return(nil)
		processor.Process("result:by-season-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing result into repository when processing result by season id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fix := []app.Fixture{*newFixture(34)}

		fixtureRepo.On("BySeasonID", uint64(34)).Return(fix, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(res, nil)
		resultRepo.On("Update", res).Return(nil)
		processor.Process("result:by-season-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert result into repository when processing result by season id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fix := []app.Fixture{*newFixture(34)}

		fixtureRepo.On("BySeasonID", uint64(34)).Return(fix, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(&app.Result{}, errors.New("not found"))
		resultRepo.On("Insert", res).Return(errors.New("error occurred"))
		processor.Process("result:by-season-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update result into repository when processing result by season id command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		fix := []app.Fixture{*newFixture(34)}

		fixtureRepo.On("BySeasonID", uint64(34)).Return(fix, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(res, nil)
		resultRepo.On("Update", res).Return(errors.New("error occurred"))
		processor.Process("result:by-season-id", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new result into repository when processing result today command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(&app.Result{}, errors.New("not found"))
		resultRepo.On("Insert", res).Return(nil)
		processor.Process("result:today", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing result into repository when processing result today command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(res, nil)
		resultRepo.On("Update", res).Return(nil)
		processor.Process("result:today", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert result into repository when processing result today command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(&app.Result{}, errors.New("not found"))
		resultRepo.On("Insert", res).Return(errors.New("error occurred"))
		processor.Process("result:today", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update result into repository when processing result today command", func(t *testing.T) {
		t.Helper()

		resultRepo := new(mock.ResultRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.ResultRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewResultProcessor(resultRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		res := newResult(34)

		results := make([]*app.Result, 1)
		results[0] = res

		ch := resultChannel(results)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("ResultsByFixtureIDs", []uint64{34}).Return(ch)
		resultRepo.On("ByFixtureID", uint64(34)).Return(res, nil)
		resultRepo.On("Update", res).Return(errors.New("error occurred"))
		processor.Process("result:today", "34", done)

		<-done

		requester.AssertExpectations(t)
		resultRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newResult(f uint64) *app.Result {
	return &app.Result{FixtureID: f}
}

func resultChannel(results []*app.Result) chan *app.Result {
	ch := make(chan *app.Result, len(results))

	for _, c := range results {
		ch <- c
	}

	close(ch)

	return ch
}