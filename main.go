package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", routePath).Methods("GET")
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func routePath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World")
	w.WriteHeader(http.StatusOK)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Healthcheck OK")
	w.WriteHeader(http.StatusOK)
}