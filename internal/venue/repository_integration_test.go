package venue

import (
	"github.com/joesweeny/statshub/internal/config"
	"database/sql"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/joesweeny/statshub/internal/model"
	"time"
)

func TestInsert(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresVenueRepository{Connection: conn}

	t.Run("increase table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newVenue(i)

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

	conn.Close()
}

func TestGetById(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresVenueRepository{Connection: conn}

	t.Run("venue can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		v := newVenue(13)

		if err := repo.Insert(v); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.GetById(13)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(13, r.ID)
		a.Equal("London Stadium", r.Name)
		a.Equal("Grass", *r.Surface)
		a.Nil(r.Address)
		a.Equal("London", *r.City)
		a.Nil(r.Capacity)
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if round does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.GetById(99)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %s", ErrNotFound, err)
		}
	})

	conn.Close()
}

func TestUpdate(t *testing.T) {
	conn, cleanUp := getConnection(t)
	repo := PostgresVenueRepository{Connection: conn}

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
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.GetById(2)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(2, r.ID)
		a.Equal("Upton Park", r.Name)
		a.Nil(r.Surface)
		a.Equal("Stratford", *r.Address)
		a.Equal("London", *r.City)
		a.Equal(60000, *r.Capacity)
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if venue does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newVenue(146)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}

		if err != ErrNotFound {
			t.Fatalf("Test failed, expected %v, got %v", ErrNotFound, err)
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
		_, err := db.Exec("delete from sportmonks_venue")
		if err != nil {
			t.Fatalf("Failed to clear database. %s", err.Error())
		}
	}
}

func newVenue(id int) *model.Venue {
	var s = "Grass"
	var c = "London"

	return &model.Venue{
		ID:      	id,
		Name:    	"London Stadium",
		Surface: 	&s,
		City: 		&c,
		CreatedAt: 	time.Unix(1548086929, 0),
		UpdatedAt: 	time.Unix(1548086929, 0),
	}
}