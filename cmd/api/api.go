package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", routePath).Methods("GET")
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func routePath(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	fmt.Fprint(w, "Hello "+name+" This actually works!!")
	w.WriteHeader(http.StatusOK)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Healthcheck OK")
	w.WriteHeader(http.StatusOK)
}
