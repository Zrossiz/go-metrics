package main

import "os"

func main() {
	defer os.Exit(1) // want "you can not use Exit in main function package main"

	defer func() {
		os.Exit(0)
	}()
}
