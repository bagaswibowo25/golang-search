package cli

import (
	"log"

	"github.com/bagaswibowo25/golang-search/pkg/app"
)

// Next:
// kubectl -> https://github.com/spf13/cobra
// release portal -> https://github.com/urfave/cli
// run --port 8080
func Run() {
	err := app.StartServer(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
