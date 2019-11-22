package postgres_test

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/postgres"
	"github.com/statistico/statistico-data/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFixtureRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_fixture")
	repo := postgres.NewFixtureRepository(conn, test.Clock)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i < 4; i++ {
			c := newFixture(uint64(i))

			if err := repo.Insert(c); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}

			row := conn.QueryRow("select count(*) from sportmonks_fixture")

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
		c := newFixture(50)

		if err := repo.Insert(c); err != nil {
			t.Errorf("Test failed, expected nil, got %s", err)
		}

		if e := repo.Insert(c); e == nil {
			t.Fatalf("Test failed, expected %s, got nil", e)
		}
	})
}

func TestFixtureRepository_ByID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_fixture")
	repo := postgres.NewFixtureRepository(conn, test.Clock)

	t.Run("fixture can be retrieved by ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newFixture(43)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		r, err := repo.ByID(43)

		if err != nil {
			t.Errorf("Error when retrieving a record from the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(43), r.ID)
		a.Equal(uint64(14567), r.SeasonID)
		a.Equal(uint64(165789), *r.RoundID)
		a.Nil(f.VenueID)
		a.Equal(uint64(451), r.HomeTeamID)
		a.Equal(uint64(924), r.AwayTeamID)
		a.Nil(r.RefereeID)
		a.Equal("2019-01-21 16:08:49 +0000 UTC", r.Date.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns error if fixture does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := repo.ByID(99)

		if err == nil {
			t.Errorf("Test failed, expected %v, got nil", err)
		}
	})
}

func TestFixtureRepository_Update(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_fixture")
	repo := postgres.NewFixtureRepository(conn, test.Clock)

	t.Run("modifies existing fixture", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		f := newFixture(78)

		if err := repo.Insert(f); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}

		var venueId = uint64(574)
		var roundId *uint64
		var d = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)

		f.VenueID = &venueId
		f.AwayTeamID = uint64(4390)
		f.RoundID = roundId
		f.Date = d

		if err := repo.Update(f); err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		r, err := repo.ByID(78)

		if err != nil {
			t.Errorf("Error when updating a record in the database: %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(78), f.ID)
		a.Equal(uint64(14567), f.SeasonID)
		a.Nil(f.RoundID)
		a.Equal(uint64(574), *f.VenueID)
		a.Equal(uint64(451), f.HomeTeamID)
		a.Equal(uint64(4390), f.AwayTeamID)
		a.Nil(f.RefereeID)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.Date.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", r.UpdatedAt.String())
	})

	t.Run("returns an error if fixture does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()
		c := newFixture(146)

		err := repo.Update(c)

		if err == nil {
			t.Fatalf("Test failed, expected nil, got %v", err)
		}
	})
}

func TestFixtureRepository_GetIDs(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_fixture")
	repo := postgres.NewFixtureRepository(conn, test.Clock)

	t.Run("returns int slice of fixture ids where date is between two dates", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		for i := 1; i <= 4; i++ {
			s := newFixture(uint64(i))

			if err := repo.Insert(s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		for i := 5; i <= 8; i++ {
			s := app.Fixture{
				ID:         uint64(i),
				SeasonID:   uint64(14567),
				HomeTeamID: uint64(451),
				AwayTeamID: uint64(924),
				Date:       time.Unix(1550066305, 0),
				CreatedAt:  time.Unix(1546965200, 0),
				UpdatedAt:  time.Unix(1546965200, 0),
			}

			if err := repo.Insert(&s); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		after := time.Unix(1548086910, 0)
		before := time.Unix(1548086950, 0)

		query := app.FixtureRepositoryQuery{DateTo: &before, DateFrom: &after}

		ids, err := repo.GetIDs(query)

		want := []uint64{1, 2, 3, 4}

		if err != nil {
			t.Fatalf("Test failed, expected %v, got %s", want, err.Error())
		}

		all, err := repo.GetIDs(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected %v, got %s", want, err.Error())
		}

		assert.Equal(t, 8, len(all))
		assert.Equal(t, want, ids)
	})
}

func TestFixtureRepository_ByTeamID(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_fixture")
	repo := postgres.NewFixtureRepository(conn, test.Clock)

	t.Run("returns slice of fixture struct matching parameters provided", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		fix, err := repo.ByTeamID(66, 100, time.Unix(1550066317, 0))

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 3, len(fix))
	})

	t.Run("results can be filtered by limit", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		fix, err := repo.ByTeamID(66, 1, time.Unix(1550066317, 0))

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 1, len(fix))
		assert.Equal(t, uint64(6), fix[0].ID)
	})

	t.Run("empty result set returned if no results match parameters", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		fix, err := repo.ByTeamID(14059, 1, time.Unix(1550066317, 0))

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 0, len(fix))
	})
}

func TestFixtureRepository_Get(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, "sportmonks_fixture")
	repo := postgres.NewFixtureRepository(conn, test.Clock)

	t.Run("returns slice of fixture struct matching with matching season ID", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		season := uint64(6012)

		query := app.FixtureRepositoryQuery{SeasonID: &season}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 4, len(fix))

		for i := 0; i <= 1; i++ {
			f := fix[i]
			assert.Equal(t, uint64(i+5), f.ID)
			assert.Equal(t, uint64(6012), f.SeasonID)
		}
	})

	t.Run("returns slice of fixture struct matching season id and date before parameters provided", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		season := uint64(6012)
		before := time.Unix(1550066319, 0)

		query := app.FixtureRepositoryQuery{SeasonID: &season, DateTo: &before}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 4, len(fix))
		assert.Equal(t, uint64(5), fix[0].ID)
		assert.Equal(t, uint64(6), fix[1].ID)
		assert.Equal(t, uint64(7), fix[2].ID)
		assert.Equal(t, uint64(8), fix[3].ID)

		for _, f := range fix {
			assert.Equal(t, uint64(6012), f.SeasonID)
		}
	})

	t.Run("returns slice of fixture struct restricted by before date", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		before := time.Unix(1550066312, 0)

		query := app.FixtureRepositoryQuery{DateTo: &before}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 5, len(fix))
	})

	t.Run("returns slice of fixture struct restricted by before and after dates", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		season := uint64(6012)
		from := time.Unix(1550066313, 0)
		to := time.Unix(1550066316, 0)

		query := app.FixtureRepositoryQuery{DateTo: &to, DateFrom: &from, SeasonID: &season}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 2, len(fix))
	})

	t.Run("returns slice of fixture struct restricted by after date", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		after := time.Unix(1550066316, 0)

		query := app.FixtureRepositoryQuery{DateFrom: &after}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 3, len(fix))
	})

	t.Run("returns slice of fixture struct matching team id parameters provided", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		home := uint64(66)
		away := uint64(924)

		query := app.FixtureRepositoryQuery{HomeTeamID: &home, AwayTeamID: &away}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 4, len(fix))
	})

	t.Run("empty result set returned if no results match parameters", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertFixtures(t, repo)

		season := uint64(999999999)

		query := app.FixtureRepositoryQuery{SeasonID: &season}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		all, err := repo.Get(app.FixtureRepositoryQuery{})

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		assert.Equal(t, 9, len(all))
		assert.Equal(t, 0, len(fix))
	})

	t.Run("returns fixture struct ordered by date asc", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		fixtures := []app.Fixture{
			{
				ID:         uint64(1),
				SeasonID:   uint64(14567),
				HomeTeamID: 451,
				AwayTeamID: 924,
				Date:       time.Unix(1548086929, 0),
			},
			{
				ID:         uint64(2),
				SeasonID:   uint64(14567),
				HomeTeamID: 451,
				AwayTeamID: 924,
				Date:       time.Unix(1548086920, 0),
			},
			{
				ID:         uint64(3),
				SeasonID:   uint64(14567),
				HomeTeamID: 451,
				AwayTeamID: 924,
				Date:       time.Unix(1548086925, 0),
			},
		}

		for _, fix := range fixtures {
			if err := repo.Insert(&fix); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		sort := "date_asc"

		query := app.FixtureRepositoryQuery{SortBy: &sort}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(uint64(2), fix[0].ID)
		a.Equal(int64(1548086920), fix[0].Date.Unix())
		a.Equal(uint64(3), fix[1].ID)
		a.Equal(int64(1548086925), fix[1].Date.Unix())
		a.Equal(uint64(1), fix[2].ID)
		a.Equal(int64(1548086929), fix[2].Date.Unix())
	})

	t.Run("returns fixture struct ordered by date desc", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		fixtures := []app.Fixture{
			{
				ID:         uint64(1),
				SeasonID:   uint64(14567),
				HomeTeamID: 451,
				AwayTeamID: 924,
				Date:       time.Unix(1548086929, 0),
			},
			{
				ID:         uint64(2),
				SeasonID:   uint64(14567),
				HomeTeamID: 451,
				AwayTeamID: 924,
				Date:       time.Unix(1548086920, 0),
			},
			{
				ID:         uint64(3),
				SeasonID:   uint64(14567),
				HomeTeamID: 451,
				AwayTeamID: 924,
				Date:       time.Unix(1548086925, 0),
			},
		}

		for _, fix := range fixtures {
			if err := repo.Insert(&fix); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		sort := "date_desc"

		query := app.FixtureRepositoryQuery{SortBy: &sort}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(uint64(2), fix[2].ID)
		a.Equal(int64(1548086920), fix[2].Date.Unix())
		a.Equal(uint64(3), fix[1].ID)
		a.Equal(int64(1548086925), fix[1].Date.Unix())
		a.Equal(uint64(1), fix[0].ID)
		a.Equal(int64(1548086929), fix[0].Date.Unix())
	})

	t.Run("returns fixture struct slice where home team name like", func(t *testing.T) {
		conn, clean := test.GetConnection(t, "sportmonks_team")
		teamRepo := postgres.NewTeamRepository(conn, test.Clock)

		t.Helper()
		defer clean()
		defer cleanUp()

		insertFixtures(t, repo)

		teams := []app.Team{
			{
				ID:           451,
				Name:         "West Ham United",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
			{
				ID:           22,
				Name:         "Chelsea",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
		}

		for _, team := range teams {
			if err := teamRepo.Insert(&team); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		home := "West"

		query := app.FixtureRepositoryQuery{HomeTeamNameLike: &home}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(4, len(fix))
	})

	t.Run("returns fixture struct slice where away team name like", func(t *testing.T) {
		conn, clean := test.GetConnection(t, "sportmonks_team")
		teamRepo := postgres.NewTeamRepository(conn, test.Clock)

		t.Helper()
		defer clean()
		defer cleanUp()

		insertFixtures(t, repo)

		teams := []app.Team{
			{
				ID:           451,
				Name:         "West Ham United",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
			{
				ID:           924,
				Name:         "Chelsea",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
		}

		for _, team := range teams {
			if err := teamRepo.Insert(&team); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		away := "Chel"

		query := app.FixtureRepositoryQuery{AwayTeamNameLike: &away}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(8, len(fix))
	})

	t.Run("returns fixture struct slice where home team and away team name like", func(t *testing.T) {
		conn, clean := test.GetConnection(t, "sportmonks_team")
		teamRepo := postgres.NewTeamRepository(conn, test.Clock)

		t.Helper()
		defer clean()
		defer cleanUp()

		insertFixtures(t, repo)

		teams := []app.Team{
			{
				ID:           451,
				Name:         "West Ham United",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
			{
				ID:           924,
				Name:         "Chelsea",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
		}

		for _, team := range teams {
			if err := teamRepo.Insert(&team); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		home := "West"
		away := "Chel"

		query := app.FixtureRepositoryQuery{HomeTeamNameLike: &home, AwayTeamNameLike: &away}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(4, len(fix))
	})

	t.Run("returns fixture struct slice where away team name like", func(t *testing.T) {
		conn, clean := test.GetConnection(t, "sportmonks_team")
		teamRepo := postgres.NewTeamRepository(conn, test.Clock)

		t.Helper()
		defer clean()
		defer cleanUp()

		insertFixtures(t, repo)

		teams := []app.Team{
			{
				ID:           451,
				Name:         "West Ham United",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
			{
				ID:           924,
				Name:         "Chelsea",
				VenueID:      560,
				CountryID:    uint64(462),
				NationalTeam: false,
			},
		}

		for _, team := range teams {
			if err := teamRepo.Insert(&team); err != nil {
				t.Errorf("Error when inserting record into the database: %s", err.Error())
			}
		}

		away := "Chel"

		query := app.FixtureRepositoryQuery{AwayTeamNameLike: &away}

		fix, err := repo.Get(query)

		if err != nil {
			t.Fatalf("Test failed, expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(8, len(fix))
	})
}

func newFixture(id uint64) *app.Fixture {
	var roundId = uint64(165789)

	return &app.Fixture{
		ID:         id,
		SeasonID:   uint64(14567),
		RoundID:    &roundId,
		HomeTeamID: 451,
		AwayTeamID: 924,
		Date:       time.Unix(1548086929, 0),
	}
}

func insertFixtures(t *testing.T, repo app.FixtureRepository) {
	for i := 1; i <= 4; i++ {
		s := newFixture(uint64(i))

		if err := repo.Insert(s); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}
	}

	for i := 5; i <= 8; i++ {
		x := 1550066310 + i
		s := app.Fixture{
			ID:         uint64(i),
			SeasonID:   uint64(6012),
			HomeTeamID: uint64(66),
			AwayTeamID: uint64(924),
			Date:       time.Unix(int64(x), 0),
			CreatedAt:  time.Unix(1546965200, 0),
			UpdatedAt:  time.Unix(1546965200, 0),
		}

		if err := repo.Insert(&s); err != nil {
			t.Errorf("Error when inserting record into the database: %s", err.Error())
		}
	}

	s := app.Fixture{
		ID:         uint64(99),
		SeasonID:   uint64(145),
		HomeTeamID: uint64(66),
		AwayTeamID: uint64(32),
		Date:       time.Unix(1550066312, 0),
		CreatedAt:  time.Unix(1546965200, 0),
		UpdatedAt:  time.Unix(1546965200, 0),
	}

	if err := repo.Insert(&s); err != nil {
		t.Errorf("Error when inserting record into the database: %s", err.Error())
	}
}
