package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/peterh/liner"
)

// Resolve asks the user to resolve the identifiers
func Resolve(identifiers []string, config *Config, in *bufio.Reader) map[string]string {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	values := make(map[string]string, len(identifiers))
	for _, identifier := range identifiers {
		if _, ok := values[identifier]; !ok {
			prompt := fmt.Sprintf("%s: ", identifier)
			var text string
			var err error
			if in == nil {
				line.ClearHistory()
				hs := config.history(id)
				for i := len(hs) - 1; i >= 0; i-- {
					line.AppendHistory(hs[i])
				}
				text, err = line.Prompt(prompt)
				if err != nil {
					if err == liner.ErrPromptAborted || err == io.EOF {
						os.Exit(1)
					}
					log.Fatal(err)
				}
			} else {
				fmt.Printf(prompt)
				text, err = in.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
			}
			values[identifier] = strings.TrimSuffix(text, "\n")
		}
	}
	return values
}
