package main

import (
	"fmt"
	"io"
	"os"
)

var (
	name        = "fillin"
	version     = "0.3.3"
	description = "fill-in your command and execute"
	author      = "itchyny"
)

func main() {
	if err := Exec(); err != nil {
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		}
		os.Exit(1)
	}
}
