package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSeasonRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_season")
	repo := postgres.NewSeasonRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			s := newSeason(int64(i), true)

			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_season")

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
		c := newSeason(10, true)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestSeasonRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_season")
	repo := postgres.NewSeasonRepository(conn, test.Clock)

	t.Run("modifies existing record", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		s := newSeason(50, true)

		if err := repo.Insert(s); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		s.IsCurrent = false
		s.CompetitionID = int64(2)
		s.UpdatedAt = d

		if err := repo.Update(s); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByID(50)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(50), r.ID)
		a.Equal("2018-2019", r.Name)
		a.Equal(int64(2), r.CompetitionID)
		a.Equal(false, r.IsCurrent)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if record does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newSeason(146, true)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func TestSeasonRepository_ByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_season")
	repo := postgres.NewSeasonRepository(conn, test.Clock)

	t.Run("season can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		s := newSeason(146, true)

		err := repo.Update(s)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}

		if err := repo.Insert(s); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByID(146)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(146), r.ID)
		a.Equal("2018-2019", r.Name)
		a.Equal(int64(560), r.CompetitionID)
		a.True(r.IsCurrent)
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-08 16:33:20 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if season does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.ByID(4); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestSeasonRepository_IDs(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_season")
	repo := postgres.NewSeasonRepository(conn, test.Clock)

	t.Run("test returns a slice of int ids", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i <= 4; i++ {
			s := newSeason(int64(i), true)

			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the dataCurrentSeasonIdsbase: %s", err.Error())
			}
		}

		ids, err := repo.IDs()

		want := []int64{1, 2, 3, 4}

		if err != nil {
			t.Fatalf("Test failed, expected %v, got %s", want, err.Error())
		}

		assert.Equal(t, want, ids)
	})
}

func TestSeasonRepository_CurrentSeasonIDs(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_season")
	repo := postgres.NewSeasonRepository(conn, test.Clock)

	t.Run("returns records with is current season set to true", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		var seasons []app.Season

		for i := 1; i <= 4; i++ {
			s := newSeason(int64(i), true)

			seasons = append(seasons, *s)

			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		if err := repo.Insert(newSeason(10, false)); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		retrieved, err := repo.CurrentSeasonIDs()

		if err != nil {
			t.Fatalf("Test failed, expected %v, got %s", seasons, err.Error())
		}

		assert.Equal(t, []int64{1, 2, 3, 4}, retrieved)
	})
}

func newSeason(id int64, current bool) *app.Season {
	return &app.Season{
		ID:        id,
		Name:      "2018-2019",
		CompetitionID:  int64(560),
		IsCurrent: current,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}

