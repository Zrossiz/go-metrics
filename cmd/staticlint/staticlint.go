package main

import "github.com/Zrossiz/go-metrics/pkg/multicheker"

func main() {
	chkr := multicheker.New()
	chkr.Run()
}
