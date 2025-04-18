package main

import (
	"log"
	"net/http"

	"richisntreal-backend/cmd/bootstrap"
	"richisntreal-backend/cmd/config"
)

func main() {
	router := bootstrap.NewRouter()

	port := config.Get().App.Port
	log.Printf("🚀 Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
