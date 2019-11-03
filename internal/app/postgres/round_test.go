package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRoundRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_round")
	repo := postgres.NewRoundRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newRound(uint64(i))

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_round")

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
		c := newRound(50)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestRoundRepository_GetByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_round")
	repo := postgres.NewRoundRepository(conn, test.Clock)

	t.Run("round can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newRound(43)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByID(43)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(43), r.ID)
		a.Equal("5", r.Name)
		a.Equal(uint64(4387), r.SeasonID)
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.StartDate.String())
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.EndDate.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if round does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.ByID(4); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestRoundRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_round")
	repo := postgres.NewRoundRepository(conn, test.Clock)

	t.Run("modifies existing round", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		r := newRound(897)

		if err := repo.Insert(r); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var s = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
		var e = time.Date(2019, 01, 14, 11, 29, 00, 00, time.UTC)

		r.StartDate = s
		r.EndDate = e

		if err := repo.Update(r); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByID(897)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(897), r.ID)
		a.Equal("5", r.Name)
		a.Equal(uint64(4387), r.SeasonID)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.StartDate.String())
		a.Equal("2019-01-14 11:29:00 +0000 UTC", r.EndDate.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if round does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newRound(146)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func newRound(id uint64) *app.Round {
	return &app.Round{
		ID:        id,
		Name:      "5",
		SeasonID:  uint64(4387),
		StartDate: time.Unix(1548086929, 0),
		EndDate:   time.Unix(1548086929, 0),
	}
}
