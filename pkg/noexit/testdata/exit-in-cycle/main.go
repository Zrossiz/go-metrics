package main

import "os"

func main() {
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			os.Exit(i) // want "you can not use Exit in main function package main"
		}
	}
}
