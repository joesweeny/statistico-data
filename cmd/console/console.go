package main

import (
	"flag"
	"fmt"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/joesweeny/statshub/internal/bootstrap"
	"os"
)

const competition = "competition"
const country = "country"
const fixture = "fixture"
const fixtureCurrentSeason = "fixture:current-season"
const season = "season"

var option = flag.String("option", "", "Provide the model name to process")

func main() {
	app := bootstrap.Bootstrap{Config: config.GetConfig()}

	flag.Parse()

	var service bootstrap.Service

	switch *option {
	case competition:
		service = app.GetCompetitionService()
		break
	case country:
		service = app.GetCountryService()
		break
	case fixture:
		service = app.GetFixtureService()
		break
	case season:
		service = app.GetSeasonService()
	default:
		fmt.Println("The option provided is not supported")
		os.Exit(1)
	}

	if err := service.Process(); err != nil {
		fail(option, err)
		os.Exit(1)
	}

	fmt.Printf("Processing complete for %s\n", *option)
	os.Exit(0)
}

func fail(model *string, err error) {
	if err != nil {
		fmt.Printf("Error when processing %s: %s\n", *model, err.Error())
		os.Exit(1)
	}
}

