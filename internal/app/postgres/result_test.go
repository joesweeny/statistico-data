package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestResultRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_result")
	repo := postgres.NewResultRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newResult(uint64(i))

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_result")

			var count int

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, i, count)
		}
	})

	t.Run("returns error when ID primary key violates unique constraint", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newResult(50)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}
	})
}

func TestResultRepository_ByFixtureID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_result")
	repo := postgres.NewResultRepository(conn, test.Clock)

	t.Run("result can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newResult(50)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByFixtureID(50)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(50), r.FixtureID)
		a.Nil(r.PitchCondition)
		a.Nil(r.HomeFormation)
		a.Nil(r.AwayFormation)
		a.Nil(r.HomeScore)
		a.Nil(r.AwayScore)
		a.Nil(r.HomePenScore)
		a.Nil(r.AwayPenScore)
		a.Nil(r.HalfTimeScore)
		a.Nil(r.FullTimeScore)
		a.Nil(r.ExtraTimeScore)
		a.Nil(r.HomeLeaguePosition)
		a.Nil(r.AwayLeaguePosition)
		a.Nil(r.Minutes)
		a.Nil(r.AddedTime)
		a.Nil(r.ExtraTime)
		a.Nil(r.InjuryTime)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if result does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.ByFixtureID(99)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestUpdate(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_result")
	repo := postgres.NewResultRepository(conn, test.Clock)

	t.Run("modifies existing result", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newResult(78)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var pitch = "Good"
		var homeFormation = "4-2-3-1"
		var awayFormation = "5-3-2"
		var homeScore = 4
		var awayScore = 1
		var awayPenScore = 1
		var halfTimeScore = "0-0"
		var fullTimeScore = "4-1"
		var homePosition = 2
		var awayPosition = 18
		var mins = 90
		var added = 5
		var injury = 2
		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		m.PitchCondition = &pitch
		m.HomeFormation = &homeFormation
		m.AwayFormation = &awayFormation
		m.HomeScore = &homeScore
		m.AwayScore = &awayScore
		m.AwayPenScore = &awayPenScore
		m.HalfTimeScore = &halfTimeScore
		m.FullTimeScore = &fullTimeScore
		m.HomeLeaguePosition = &homePosition
		m.AwayLeaguePosition = &awayPosition
		m.Minutes = &mins
		m.AddedTime = &added
		m.InjuryTime = &injury
		m.UpdatedAt = d

		if err := repo.Update(m); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByFixtureID(78)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(78), r.FixtureID)
		a.Equal("Good", *r.PitchCondition)
		a.Equal("4-2-3-1", *r.HomeFormation)
		a.Equal("5-3-2", *r.AwayFormation)
		a.Equal(4, *r.HomeScore)
		a.Equal(1, *r.AwayScore)
		a.Nil(r.HomePenScore)
		a.Equal(1, *r.AwayPenScore)
		a.Equal("0-0", *r.HalfTimeScore)
		a.Equal("4-1", *r.FullTimeScore)
		a.Nil(r.ExtraTimeScore)
		a.Equal(2, *r.HomeLeaguePosition)
		a.Equal(18, *r.AwayLeaguePosition)
		a.Equal(90, *r.Minutes)
		a.Equal(5, *r.AddedTime)
		a.Nil(r.ExtraTime)
		a.Equal(2, *r.InjuryTime)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if result does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		r := newResult(146)

		err := repo.Update(r)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func newResult(f uint64) *app.Result {
	return &app.Result{FixtureID: f}
}
