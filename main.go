package main

import (
	"log"
)

var name = "fillin"
var version = "v0.0.3"
var description = "fill-in your command line"
var author = "itchyny"

func main() {
	if err := Exec(); err != nil {
		log.Fatal(err)
	}
}
