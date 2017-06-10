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
func Resolve(identifiers []*Identifier, config *Config, in *bufio.Reader, out *bufio.Writer) map[string]map[string]string {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	values := make(map[string]map[string]string)
	for _, id := range identifiers {
		if !found(values, id) {
			prompt := fmt.Sprintf("%s: ", id.key)
			if id.scope != "" {
				prompt = fmt.Sprintf("[%s] %s: ", id.scope, id.key)
			}
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
				out.WriteString(prompt)
				out.Flush()
				text, err = in.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
			}
			insert(values, id, strings.TrimSuffix(text, "\n"))
		}
	}
	return values
}
