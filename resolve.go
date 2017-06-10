package main

import (
	"bufio"
	"fmt"
	"strings"
)

// Resolve asks the user to resolve the identifiers
func Resolve(identifiers []string, config *Config, in *bufio.Reader) map[string]string {
	values := make(map[string]string, len(identifiers))
	for _, identifier := range identifiers {
		if _, ok := values[identifier]; !ok {
			fmt.Printf("%s: ", identifier)
			text, _ := in.ReadString('\n')
			values[identifier] = strings.TrimSuffix(text, "\n")
		}
	}
	return values
}
