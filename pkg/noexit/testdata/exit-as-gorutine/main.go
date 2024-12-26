package main

import "os"

func main() {
	go os.Exit(1) // want "you can not use Exit in main function package main"
}
