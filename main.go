package main

import (
	"log"
)

var name = "fillin"
var version = "v0.0.4"
var description = "fill-in your command and execute"
var author = "itchyny"

func main() {
	if err := Exec(); err != nil {
		log.Fatal(err)
	}
}
