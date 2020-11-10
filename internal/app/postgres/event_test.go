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
			m := newGoalEvent(uint64(i), 45)

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
		m := newGoalEvent(50, 45)

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

		m := newGoalEvent(33, 45)

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

func TestEventRepository_InsertCardEvent(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_card_event")
	repo := postgres.NewEventRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newCardEvent(uint64(i), 45)

			if err := repo.InsertCardEvent(m); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_card_event")

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
		m := newCardEvent(50, 45)

		if err := repo.InsertCardEvent(m); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.InsertCardEvent(m); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestEventRepository_CardEventsForFixture(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_card_event")
	repo := postgres.NewEventRepository(conn, test.Clock)

	t.Run("returns a slice of card event struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		events := []*app.CardEvent{
			newCardEvent(1, 45),
			newCardEvent(2, 102),
			newCardEvent(3, 45),
			newCardEvent(4, 45),
		}

		for _, e := range events {
			if err := repo.InsertCardEvent(e); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		fetched, err := repo.CardEventsForFixture(45)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err)
		}

		assert.Equal(t, 3, len(fetched))

		for _, e := range fetched {
			assert.Equal(t, uint64(45), e.FixtureID)
		}
	})
}

func TestEventRepository_GoalEventsForFixture(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_goal_event")
	repo := postgres.NewEventRepository(conn, test.Clock)

	t.Run("returns a slice of goal event struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		events := []*app.GoalEvent{
			newGoalEvent(1, 45),
			newGoalEvent(2, 102),
			newGoalEvent(3, 45),
			newGoalEvent(4, 45),
		}

		for _, e := range events {
			if err := repo.InsertGoalEvent(e); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		fetched, err := repo.GoalEventsForFixture(45)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err)
		}

		assert.Equal(t, 3, len(fetched))

		for _, e := range fetched {
			assert.Equal(t, uint64(45), e.FixtureID)
		}
	})
}

func newCardEvent(id, fixtureID uint64) *app.CardEvent {
	return &app.CardEvent{
		ID:          id,
		TeamID:      uint64(4509),
		Type:        "red",
		FixtureID:   fixtureID,
		PlayerID:    uint64(3401),
		Minute:      85,
		Reason:      nil,
		CreatedAt:   time.Time{},
	}
}

func newGoalEvent(id, fixtureID uint64) *app.GoalEvent {
	return &app.GoalEvent{
		ID:        id,
		FixtureID: fixtureID,
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
