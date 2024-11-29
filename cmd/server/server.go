// Package main is the entry point of the application. It initializes and starts the server.
package main

import (
	"fmt"

	"github.com/Zrossiz/go-metrics/internal/server/app"
)

// main is the entry point of the application.
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	app.StartServer()
}
