package main

import (
	cli "github.com/bagaswibowo25/golang-search/pkg/cli"
)

// Config: config, terima configuration dari end users. Via CLI, via file .ini, atau env vars
// App: Handle business logic, yg specific ke use case kita
// Server: Handle / interface network (HTTP)
// CLI: Start this software via CLI
//
// User -> *CLI -> *Config + HTTP Server -> App -> Index, NATS, etc
//
//	--help
//	--...
func main() {
	cli.Run()
}
