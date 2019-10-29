package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPlayerRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_player")
	repo := postgres.NewPlayerRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newPlayer(int64(i))

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_player")

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
		c := newPlayer(50)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestPlayerRepository_ByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_player")
	repo := postgres.NewPlayerRepository(conn, test.Clock)

	t.Run("player can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newPlayer(int64(43))

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByID(43)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(int64(43), r.ID)
		a.Equal(int64(154), m.CountryId)
		a.Equal("Manuel", m.FirstName)
		a.Equal("Lanzini", m.LastName)
		a.Equal("Buenos Aires", *m.BirthPlace)
		a.Equal("1984-03-12", *m.DateOfBirth)
		a.Equal(3, m.PositionID)
		a.Equal("path/to/image", m.Image)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if player does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.ByID(4); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestPlayerRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_player")
	repo := postgres.NewPlayerRepository(conn, test.Clock)

	t.Run("modifies existing player", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		m := newPlayer(46)

		if err := repo.Insert(m); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		m.UpdatedAt = d

		if err := repo.Update(m); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByID(46)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(int64(46), r.ID)
		a.Equal(int64(154), m.CountryId)
		a.Equal("Manuel", m.FirstName)
		a.Equal("Lanzini", m.LastName)
		a.Equal("Buenos Aires", *m.BirthPlace)
		a.Equal("1984-03-12", *m.DateOfBirth)
		a.Equal(3, m.PositionID)
		a.Equal("path/to/image", m.Image)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if player does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newPlayer(146)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func newPlayer(id int64) *app.Player {
	var place = "Buenos Aires"
	var dob = "1984-03-12"
	return &app.Player{
		ID:          id,
		CountryId:   int64(154),
		FirstName:   "Manuel",
		LastName:    "Lanzini",
		BirthPlace:  &place,
		DateOfBirth: &dob,
		PositionID:  3,
		Image:       "path/to/image",
		CreatedAt:   time.Unix(1546965200, 0),
		UpdatedAt:   time.Unix(1546965200, 0),
	}
}
