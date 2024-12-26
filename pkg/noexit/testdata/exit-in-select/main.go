package main

import (
	"os"
	"time"
)

func main() {
	signal := make(chan os.Signal)
	ticker := time.NewTicker(100)

	defer ticker.Stop()

	select {
	case <-signal:
		os.Exit(0) // want "you can not use Exit in main function package main"

	case <-ticker.C:
		os.Exit(1) // want "you can not use Exit in main function package main"
	}
}
