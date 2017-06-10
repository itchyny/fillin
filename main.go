package main

import (
	"log"
	"os"
)

func main() {
	if err := Run(os.Args[1:], os.Environ()); err != nil {
		log.Fatal(err)
	}
}
