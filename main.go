package main

import (
	"log"
)

func main() {
	if err := Exec(); err != nil {
		log.Fatal(err)
	}
}
