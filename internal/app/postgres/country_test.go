package postgres_test

import (
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/statistico/statistico-football-data/internal/app/postgres"
	"github.com/statistico/statistico-football-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountryRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_country")
	repo := postgres.NewCountryRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newCountry(uint64(i))

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_country")

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
		c := newCountry(10)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestCountryRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_country")
	repo := postgres.NewCountryRepository(conn, test.Clock)

	t.Run("modifies existing record", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newCountry(100)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		s := "Scotland"
		c.Name = s

		if err := repo.Update(c); err != nil {
			t.Fatalf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByID(c.ID)

		if err != nil {
			t.Fatalf("Error when updating a record in the database: %s", err.Error())
		}

		got := r.Name
		want := s

		assert.Equal(t, got, want)
	})

	t.Run("returns error if record does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newCountry(146)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func TestCountryRepository_GetById(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_country")
	repo := postgres.NewCountryRepository(conn, test.Clock)

	t.Run("country can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		c := newCountry(62)

		if err := repo.Insert(c); err != nil {
			t.Fatalf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByID(62)

		if err != nil {
			t.Fatalf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(62), r.ID)
		a.Equal("England", r.Name)
		a.Equal("Europe", r.Continent)
		a.Equal("ENG", r.ISO)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if country does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.ByID(4); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})
}

func newCountry(id uint64) *app.Country {
	c := app.Country{
		ID:        id,
		Name:      "England",
		Continent: "Europe",
		ISO:       "ENG",
	}

	return &c
}
