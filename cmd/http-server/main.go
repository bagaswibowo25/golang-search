package main

import (
	"log"

	"github.com/bagaswibowo25/golang-search/pkg/app"
)

func main() {

	// Start the HTTP server with the configured number of workers
	err := app.StartServer(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
