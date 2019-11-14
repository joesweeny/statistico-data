package main

import (
	"github.com/statistico/statistico-data/internal/app/rest"
	"log"
	"net/http"
)

func main() {
	router := rest.Router()

	log.Fatal(http.ListenAndServe(":80", router))
}
