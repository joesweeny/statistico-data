package main

import (
	"fmt"
	"github.com/statistico/statistico-data/internal/bootstrap"
	"os"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())

	statsRepo := app.TeamStatsRepository()
	resultRepo := app.ResultRepository()
	fixtureRepo := app.FixtureRepository()

	stats, err := statsRepo.Get()

	if err != nil {
		fmt.Printf("Error fetching stats from repository, %s", err.Error())
		os.Exit(1)
	}

	fmt.Print("Starting migration")

	for _, stat := range stats {
		result, err := resultRepo.ByFixtureID(stat.FixtureID)

		if err != nil {
			fmt.Printf("Error fetching Result %d\n", stat.FixtureID)
			continue
		}

		fixture, err := fixtureRepo.ByID(stat.FixtureID)

		if err != nil {
			fmt.Printf("Error fetching Fixture %d\n", stat.FixtureID)
			continue
		}

		if fixture.HomeTeamID == stat.TeamID {
			stat.Goals = result.HomeScore
		}

		if fixture.AwayTeamID == stat.TeamID {
			stat.Goals = result.AwayScore
		}

		if err := statsRepo.UpdateTeamStats(stat); err != nil {
			fmt.Printf("Error updating stat in team stats repository %d\n", stat.FixtureID)
		}

		fmt.Printf("Goals migrated for Fixture %d\n", stat.FixtureID)
	}

	fmt.Print("Migration complete")
}