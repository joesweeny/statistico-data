package process_test

import (
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/process"
	"testing"
)

func TestSeasonProcessor_Process(t *testing.T) {
	t.Run("inserts new season", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		requester := new(mock.SeasonRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSeasonProcessor(repo, requester, logger)
	})
}
