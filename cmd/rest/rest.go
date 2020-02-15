package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/statistico/statistico-data/internal/app/rest"
	"github.com/statistico/statistico-data/internal/bootstrap"
	"log"
	"net/http"
)

func main() {
	container := bootstrap.BuildContainer(bootstrap.BuildConfig())

	router := httprouter.New()
	router.GET("/", rest.RoutePath)
	router.GET("/healthcheck", rest.HealthCheck)
	router.GET("/openapi.json", rest.RenderApiDocs)
	router.GET("/season/:id/fixtures", container.RestFixtureHandler().SeasonFixtures)

	log.Fatal(http.ListenAndServe(":80", router))
}
