package main

import (
	"log"
	"os"
)

func main() {
	configPath := "~/.config/fillin/fillin.json"
	if err := Run(configPath, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
