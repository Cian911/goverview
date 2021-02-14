package main

import (
	"log"
	"net/http"

	"github.com/cian911/goverview/api"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	api.HandleRoutes(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
