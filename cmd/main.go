package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/cian911/goverview/api"
	"github.com/cian911/goverview/pkg/websocket"
	"github.com/gorilla/mux"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	router := mux.NewRouter().StrictSlash(true)
	api.HandleRoutes(router)

	fs := http.FileServer(http.FS(staticFiles))
	router.PathPrefix("/static/").Handler(fs)

	router.HandleFunc("/ws", socket)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func socket(w http.ResponseWriter, r *http.Request) {
	vars := map[string]string{
		"id":   r.FormValue("id"),
		"repo": r.FormValue("repo"),
	}
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	go websocket.Writer(ws, vars)
}
