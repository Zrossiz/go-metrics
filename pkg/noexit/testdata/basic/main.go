package main

import "os"

func main() {
	code := 0

	if code == 0 {
		os.Exit(0) // want "you can not use Exit in main function package main"
	}

	os.Exit(code) // want "you can not use Exit in main function package main"
}
