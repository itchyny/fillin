package main

import (
	"log"
	"os"
)

func main() {
	if err := Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
