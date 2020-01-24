package main

import (
	"fmt"
	"io"
	"os"
)

var name = "fillin"
var version = "0.3.0"
var description = "fill-in your command and execute"
var author = "itchyny"

func main() {
	if err := Exec(); err != nil {
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		}
		os.Exit(1)
	}
}
