package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVenueRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_venue")
	repo := postgres.NewVenueRepository(conn, test.Clock)

	t.Run("increase table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newVenue(int64(i))

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_venue")

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
		c := newVenue(50)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestVenueRepository_GetById(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_venue")
	repo := postgres.NewVenueRepository(conn, test.Clock)

	t.Run("venue can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		v := newVenue(13)

		if err := repo.Insert(v); err != nil {
			t.Fatalf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.GetById(13)

		if err != nil {
			t.Fatalf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(13), r.ID)
		a.Equal("London Stadium", r.Name)
		a.Equal("Grass", *r.Surface)
		a.Nil(r.Address)
		a.Equal("London", *r.City)
		a.Nil(r.Capacity)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if round does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.GetById(99); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestVenueRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_venue")
	repo := postgres.NewVenueRepository(conn, test.Clock)

	t.Run("modifies existing venue", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		v := newVenue(2)

		if err := repo.Insert(v); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var add = "Stratford"
		var c = 60000

		v.Address = &add
		v.Capacity = &c
		v.Name = "Upton Park"
		v.Surface = nil

		if err := repo.Update(v); err != nil {
			t.Fatalf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.GetById(2)

		if err != nil {
			t.Fatalf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(int64(2), r.ID)
		a.Equal("Upton Park", r.Name)
		a.Nil(r.Surface)
		a.Equal("Stratford", *r.Address)
		a.Equal("London", *r.City)
		a.Equal(60000, *r.Capacity)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if venue does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newVenue(146)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func newVenue(id int64) *app.Venue {
	var s = "Grass"
	var c = "London"

	return &app.Venue{
		ID:        id,
		Name:      "London Stadium",
		Surface:   &s,
		City:      &c,
		CreatedAt: time.Unix(1548086929, 0),
		UpdatedAt: time.Unix(1548086929, 0),
	}
}
