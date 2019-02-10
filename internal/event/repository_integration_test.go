package event

import (
	"database/sql"
	"fmt"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInsertGoalEvent(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresEventRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newGoalEvent(i)

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

	conn.Close()
}

func TestInsertSubstitutionEvent(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresEventRepository{Connection: conn}

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			m := newSubstitutionEvent(i)

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

	conn.Close()
}

func TestGoalEventById(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresEventRepository{Connection: conn}

	t.Run("goal event can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newGoalEvent(33)

		if err := repo.InsertGoalEvent(m); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		r, err := repo.GoalEventById(33)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(33, r.ID)
		a.Equal(45, r.FixtureID)
		a.Equal(4509, r.TeamID)
		a.Equal(3401, r.PlayerID)
		a.Nil(r.PlayerAssistID)
		a.Equal(82, r.Minute)
		a.Equal("0-1", r.Score)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
	})

	t.Run("returns error if goal event does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.GoalEventById(99)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %s", ErrNotFound, err)
		}
	})

	conn.Close()
}

func TestSubstitutionEventById(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresEventRepository{Connection: conn}

	t.Run("substitution event can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newSubstitutionEvent(33)

		if err := repo.InsertSubstitutionEvent(m); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		r, err := repo.SubstitutionEventById(33)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(33, r.ID)
		a.Equal(45, r.FixtureID)
		a.Equal(4509, r.TeamID)
		a.Equal(3401, r.PlayerInID)
		a.Equal(901, r.PlayerOutID)
		a.Equal(82, r.Minute)
		a.True(*r.Injured)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
	})

	t.Run("returns error if substitution event does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.SubstitutionEventById(99)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %s", ErrNotFound, err)
		}
	})

	conn.Close()
}

var db = config.GetConfig().Database

func getConnection(t *testing.T) (*sql.DB, func()) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

	db, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return db, func() {
		if _, err := db.Exec("delete from sportmonks_goal_event"); err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}

		if _, err = db.Exec("delete from sportmonks_substitution_event"); err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newGoalEvent(id int) *model.GoalEvent {
	return &model.GoalEvent{
		ID:        id,
		FixtureID: 45,
		TeamID:    4509,
		PlayerID:  3401,
		Minute:    82,
		Score:     "0-1",
		CreatedAt: time.Unix(1546965200, 0),
	}
}

func newSubstitutionEvent(id int) *model.SubstitutionEvent {
	true := true
	return &model.SubstitutionEvent{
		ID:          id,
		FixtureID: 45,
		TeamID:      4509,
		PlayerInID:  3401,
		PlayerOutID: 901,
		Minute:      82,
		Injured:     &true,
		CreatedAt:   time.Unix(1546965200, 0),
	}
}
