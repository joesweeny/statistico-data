package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventRepository_InsertGoalEvent(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_goal_event")
	repo := postgres.NewEventRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newGoalEvent(uint64(i))

			if err := repo.InsertGoalEvent(m); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_goal_event")

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
		m := newGoalEvent(50)

		if err := repo.InsertGoalEvent(m); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.InsertGoalEvent(m); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestEventRepository_InsertSubstitutionEvent(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_substitution_event")
	repo := postgres.NewEventRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newSubstitutionEvent(uint64(i))

			if err := repo.InsertSubstitutionEvent(m); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_substitution_event")

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
		m := newSubstitutionEvent(50)

		if err := repo.InsertSubstitutionEvent(m); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.InsertSubstitutionEvent(m); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestEventRepository_GoalEventByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_goal_event")
	repo := postgres.NewEventRepository(conn, test.Clock)

	t.Run("goal event can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newGoalEvent(33)

		if err := repo.InsertGoalEvent(m); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		r, err := repo.GoalEventByID(33)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(uint64(33), r.ID)
		a.Equal(uint64(45), r.FixtureID)
		a.Equal(uint64(4509), r.TeamID)
		a.Equal(uint64(3401), r.PlayerID)
		a.Nil(r.PlayerAssistID)
		a.Equal(82, r.Minute)
		a.Equal("0-1", r.Score)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
	})

	t.Run("returns error if goal event does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.GoalEventByID(99)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestEventRepository_SubstitutionEventByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_substitution_event")
	repo := postgres.NewEventRepository(conn, test.Clock)

	t.Run("substitution event can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newSubstitutionEvent(33)

		if err := repo.InsertSubstitutionEvent(m); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		r, err := repo.SubstitutionEventByID(33)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(uint64(33), r.ID)
		a.Equal(uint64(45), r.FixtureID)
		a.Equal(uint64(4509), r.TeamID)
		a.Equal(uint64(3401), r.PlayerInID)
		a.Equal(uint64(901), r.PlayerOutID)
		a.Equal(82, r.Minute)
		a.True(*r.Injured)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
	})

	t.Run("returns error if substitution event does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.SubstitutionEventByID(99)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}
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
