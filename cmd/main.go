package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv" // ← new
	"richisntreal-backend/cmd/bootstrap"
	"richisntreal-backend/cmd/config"
)

func main() {
	// 1) Load .env if present (no harm if it’s missing)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to existing environment")
	}

	// 2) Now load config via Viper (env‑vars + config file + defaults)
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 3) Bootstrap and start
	router := bootstrap.NewRouter()
	port := config.Get().App.Port
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
