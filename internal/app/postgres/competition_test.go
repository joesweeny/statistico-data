package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCompetitionRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_competition")
	repo := postgres.NewCompetitionRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newCompetition(uint64(i), 462, false)

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_competition")

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
		c := newCompetition(50, 462, false)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestCompetitionRepository_ByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_competition")
	repo := postgres.NewCompetitionRepository(conn, test.Clock)

	t.Run("competition can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		c := newCompetition(45, 462, false)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByID(45)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(45), r.ID)
		a.Equal("Premier League", r.Name)
		a.Equal(uint64(462), r.CountryID)
		a.Equal(false, r.IsCup)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if country does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.ByID(4); err == nil {
			t.Fatalf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestCompetitionRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_competition")
	repo := postgres.NewCompetitionRepository(conn, test.Clock)

	t.Run("modifies existing record", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		c := newCompetition(45, 462, false)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		c.Name = "New League Name"
		c.IsCup = true

		if err := repo.Update(c); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByID(45)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(45), r.ID)
		a.Equal("New League Name", r.Name)
		a.Equal(uint64(462), r.CountryID)
		a.Equal(true, r.IsCup)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if record does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newCompetition(146, 462, false)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func TestCompetitionRepository_Get(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_competition")
	repo := postgres.NewCompetitionRepository(conn, test.Clock)

	t.Run("returns competitions filtered by is cup filter is true", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		competitions := []*app.Competition{
			newCompetition(146, 462, false),
			newCompetition(147, 462, true),
			newCompetition(148, 462, false),
		}

		for _, comp := range competitions {
			if err := repo.Insert(comp); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		isCup := true

		filter := app.CompetitionFilterQuery{
			IsCup:      &isCup,
		}

		fetched, err := repo.Get(filter)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		assert.Equal(t, 1, len(fetched))
		assert.Equal(t, uint64(147), fetched[0].ID)
	})

	t.Run("returns competitions filtered by is cup filter is false", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		competitions := []*app.Competition{
			newCompetition(146, 462, false),
			newCompetition(147, 462, true),
			newCompetition(148, 462, false),
		}

		for _, comp := range competitions {
			if err := repo.Insert(comp); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		isCup := false

		filter := app.CompetitionFilterQuery{
			IsCup:      &isCup,
		}

		fetched, err := repo.Get(filter)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		assert.Equal(t, 2, len(fetched))
		assert.Equal(t, uint64(146), fetched[0].ID)
		assert.Equal(t, uint64(148), fetched[1].ID)
	})

	t.Run("returns all records if filter query parameters are nil", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		competitions := []*app.Competition{
			newCompetition(146, 462, false),
			newCompetition(147, 462, true),
			newCompetition(148, 462, false),
		}

		for _, comp := range competitions {
			if err := repo.Insert(comp); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		fetched, err := repo.Get(app.CompetitionFilterQuery{})

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		assert.Equal(t, 3, len(fetched))
	})

	t.Run("returns competitions filtered by country IDs", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		competitions := []*app.Competition{
			newCompetition(146, 462, false),
			newCompetition(147, 463, true),
			newCompetition(148, 464, false),
		}

		for _, comp := range competitions {
			if err := repo.Insert(comp); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		filter := app.CompetitionFilterQuery{
			CountryIds:      []uint64{462, 464},
		}

		fetched, err := repo.Get(filter)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		assert.Equal(t, 2, len(fetched))
		assert.Equal(t, uint64(146), fetched[0].ID)
		assert.Equal(t, uint64(148), fetched[1].ID)
	})
}

func newCompetition(id uint64, country uint64, isCup bool) *app.Competition {
	return &app.Competition{
		ID:        id,
		Name:      "Premier League",
		CountryID: country,
		IsCup:     isCup,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}
