package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFixtureTeamXGRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "understat_fixture_team_xg")
	repo := postgres.NewFixtureTeamXGRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			f := newFixtureTeamXG(uint64(i), uint64(i + 1))

			if err := repo.Insert(f); err != nil {
				t.Errorf("Test failed, expected nil, got %s", err)
			}

			row := conn.QueryRow("select count(*) from understat_fixture_team_xg")

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

		f := newFixtureTeamXG(25, 143)

		_ = repo.Insert(f)

		if err := repo.Insert(f); err == nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}
	})
}

func TestFixtureTeamXGRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "understat_fixture_team_xg")
	repo := postgres.NewFixtureTeamXGRepository(conn, test.Clock)

	t.Run("modifies existing resource", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newFixtureTeamXG(25, 143)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		home := float32(2.30)
		away := float32(0.35)

		f.Home = &home
		f.Away = &away

		if err := repo.Update(f); err != nil {
			t.Errorf("Error updating record to the database: %s", err.Error())
		}

		fetched, err := repo.ByID(f.ID)

		if err != nil {
			t.Errorf("Error retrieving record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(25), fetched.ID)
		a.Equal(uint64(143), fetched.FixtureID)
		a.Equal(float32(2.30), *fetched.Home)
		a.Equal(float32(0.35), *fetched.Away)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", fetched.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", fetched.UpdatedAt.String())
	})

	t.Run("returns error if updating a resource that does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newFixtureTeamXG(25, 143)

		if err := repo.Update(f); err == nil {
			t.Errorf("Test failed, expected error got nil")
		}
	})
}

func TestFixtureTeamXGRepository_ByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "understat_fixture_team_xg")
	repo := postgres.NewFixtureTeamXGRepository(conn, test.Clock)

	t.Run("fixture team xg can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newFixtureTeamXG(34, 561)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		fetched, err := repo.ByID(f.ID)

		if err != nil {
			t.Errorf("Error retrieving record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(34), fetched.ID)
		a.Equal(uint64(561), fetched.FixtureID)
		a.Equal(float32(2.50), *fetched.Home)
		a.Equal(float32(0.34), *fetched.Away)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", fetched.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", fetched.UpdatedAt.String())
	})

	t.Run("returns error if fixture team xg does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.ByID(56); err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestFixtureTeamXGRepository_ByFixtureID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "understat_fixture_team_xg")
	repo := postgres.NewFixtureTeamXGRepository(conn, test.Clock)

	t.Run("fixture team xg can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newFixtureTeamXG(34, 561)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		fetched, err := repo.ByFixtureID(f.FixtureID)

		if err != nil {
			t.Errorf("Error retrieving record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(34), fetched.ID)
		a.Equal(uint64(561), fetched.FixtureID)
		a.Equal(float32(2.50), *fetched.Home)
		a.Equal(float32(0.34), *fetched.Away)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", fetched.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", fetched.UpdatedAt.String())
	})

	t.Run("returns error if fixture team xg does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		if _, err := repo.ByID(56); err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}
	})
}

func newFixtureTeamXG(id uint64, fixID uint64) *app.FixtureTeamXG {
	h := float32(2.50)
	a := float32(0.34)
	return &app.FixtureTeamXG{
		ID:        id,
		FixtureID: fixID,
		Home:      &h,
		Away:      &a,
	}
}
