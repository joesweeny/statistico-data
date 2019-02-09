package result

import (
	"time"
	"github.com/jonboulle/clockwork"
	"github.com/joesweeny/sportmonks-go-client"
	"testing"
	"github.com/stretchr/testify/assert"
)

var t = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
var clock = clockwork.NewFakeClockAt(t)
var f = Factory{clock}

func TestCreateResult(t *testing.T) {
	t.Run("a new domain result struct is hydrated", func(t *testing.T) {
		t.Helper()

		m := f.createResult(newClientFixture(5))

		a := assert.New(t)

		a.Equal(5, m.FixtureID)
		a.Equal("Good", *m.PitchCondition)
		a.Equal("4-4-2", *m.HomeFormation)
		a.Equal("4-3-2-1", *m.AwayFormation)
		a.Equal(2, *m.HomeScore)
		a.Equal(0, *m.AwayScore)
		a.Equal(0, *m.HomePenScore)
		a.Equal(0, *m.AwayPenScore)
		a.Equal("2-0", *m.HalfTimeScore)
		a.Equal("2-0", *m.FullTimeScore)
		a.Nil(m.ExtraTimeScore)
		a.Equal(1, *m.HomeLeaguePosition)
		a.Equal(5, *m.AwayLeaguePosition)
		a.Equal(90, *m.Minutes)
		a.Equal(0, *m.Seconds)
		a.Equal(0, *m.AddedTime)
		a.Equal(0, *m.ExtraTime)
		a.Equal(0, *m.InjuryTime)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
	})
}

func newClientFixture(id int) *sportmonks.Fixture {
	zero := 0
	score := "2-0"
	pitch := "Good"
	return &sportmonks.Fixture{
		ID: id,
		LeagueID: 25,
		SeasonID: 590,
		LocalTeamID: 1,
		VisitorTeamID: 9,
		Pitch: &pitch,
		Formations: sportmonks.Formations{
			LocalteamFormation: "4-4-2",
			VisitorteamFormation: "4-3-2-1",
		},
		Scores: sportmonks.Scores{
			LocalteamScore: 2,
			VisitorteamScore: 0,
			LocalteamPenScore: &zero,
			VisitorteamPenScore: &zero,
			HtScore: &score,
			FtScore: &score,
		},
		Time: sportmonks.FixtureTime{
			Minute: 90,
			Second: &zero,
			AddedTime: &zero,
			ExtraMinute: &zero,
			InjuryTime: &zero,
		},
		Standings: sportmonks.Standings{
			LocalteamPosition: 1,
			VisitorteamPosition: 5,
		},
	}
}
