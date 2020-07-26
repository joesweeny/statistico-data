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
			s := newSeason(uint64(i), 560,"2018-2019",true)

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
		c := newSeason(10, 560,"2018-2019",true)

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

		s := newSeason(50, 560,"2018-2019",true)

		if err := repo.Insert(s); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		s.IsCurrent = false
		s.CompetitionID = uint64(2)
		s.UpdatedAt = d

		if err := repo.Update(s); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByID(50)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(50), r.ID)
		a.Equal("2018-2019", r.Name)
		a.Equal(uint64(2), r.CompetitionID)
		a.Equal(false, r.IsCurrent)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if record does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newSeason(146, 560,"2018-2019",true)

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

		s := newSeason(146, 560,"2018-2019",true)

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

		a.Equal(uint64(146), r.ID)
		a.Equal("2018-2019", r.Name)
		a.Equal(uint64(560), r.CompetitionID)
		a.True(r.IsCurrent)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
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
			s := newSeason(uint64(i), 560,"2018-2019",true)

			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the dataCurrentSeasonIdsbase: %s", err.Error())
			}
		}

		ids, err := repo.IDs()

		want := []uint64{1, 2, 3, 4}

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
			s := newSeason(uint64(i), 560,"2018-2019",true)

			seasons = append(seasons, *s)

			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		if err := repo.Insert(newSeason(10, 560,"2018-2019",false)); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		retrieved, err := repo.CurrentSeasonIDs()

		if err != nil {
			t.Fatalf("Test failed, expected %v, got %s", seasons, err.Error())
		}

		assert.Equal(t, []uint64{1, 2, 3, 4}, retrieved)
	})
}

func TestSeasonRepository_ByCompetitionId(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_season")
	repo := postgres.NewSeasonRepository(conn, test.Clock)

	t.Run("returns a slice of season struct associated to a competition", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		seasons := []*app.Season{
			newSeason(1, 16036, "2018-2019",false),
			newSeason(2, 12068, "2018-2019",false),
			newSeason(3, 16036, "2018-2019",true),
		}

		for _, s := range seasons {
			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		fetched, err := repo.ByCompetitionId(16036, "name_asc")

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 2, len(fetched))
		assert.Equal(t, uint64(1), fetched[0].ID)
		assert.Equal(t, uint64(3), fetched[1].ID)
	})

	t.Run("returned results can be sorted by name descending", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		seasons := []*app.Season{
			newSeason(1, 16036, "2019-2020",false),
			newSeason(2, 12068, "2018-2019",false),
			newSeason(3, 16036, "2018-2019",true),
			newSeason(4, 16036, "2020-2021",true),
			newSeason(5, 12068, "2018-2019",true),
		}

		for _, s := range seasons {
			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		fetched, err := repo.ByCompetitionId(16036, "name_desc")

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 3, len(fetched))
		assert.Equal(t, uint64(4), fetched[0].ID)
		assert.Equal(t, "2020-2021", fetched[0].Name)
		assert.Equal(t, uint64(1), fetched[1].ID)
		assert.Equal(t, "2019-2020", fetched[1].Name)
		assert.Equal(t, uint64(3), fetched[2].ID)
		assert.Equal(t, "2018-2019", fetched[2].Name)
	})
}

func newSeason(id uint64, competitionId uint64, name string, current bool) *app.Season {
	return &app.Season{
		ID:            id,
		Name:          name,
		CompetitionID: competitionId,
		IsCurrent:     current,
		CreatedAt:     time.Unix(1546965200, 0),
		UpdatedAt:     time.Unix(1546965200, 0),
	}
}
