package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Resolve asks the user to resolve the identifiers
func Resolve(identifiers []string, config *Config) map[string]string {
	reader := bufio.NewReader(os.Stdin)
	values := make(map[string]string, len(identifiers))
	for _, identifier := range identifiers {
		if _, ok := values[identifier]; !ok {
			fmt.Printf("%s: ", identifier)
			text, _ := reader.ReadString('\n')
			values[identifier] = strings.TrimSuffix(text, "\n")
		}
	}
	return values
}
