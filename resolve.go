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

	scopeAsked := make(map[string]bool)
	if in == nil {
		for _, id := range identifiers {
			if !found(values, id) && id.scope != "" && !scopeAsked[id.scope] {
				scopeAsked[id.scope] = true
				var keys []string
				added := make(map[string]bool)
				for _, d := range identifiers {
					if id.scope == d.scope && !added[d.key] {
						keys = append(keys, d.key)
						added[d.key] = true
					}
				}
				if len(keys) > 0 {
					hs := config.historyPairs(&IdentifierGroup{scope: id.scope, keys: keys})
					if len(hs) == 0 {
						continue
					}
					line.ClearHistory()
					for i := len(hs) - 1; i >= 0; i-- {
						line.AppendHistory(hs[i])
					}
					prompt := fmt.Sprintf("[%s] %s: ", id.scope, strings.Join(keys, ", "))
					text, err := line.Prompt(prompt)
					if err != nil {
						if err == liner.ErrPromptAborted || err == io.EOF {
							os.Exit(1)
						}
						log.Fatal(err)
					}
					xs := strings.Split(strings.TrimSuffix(text, "\n"), ", ")
					if len(xs) == len(keys) {
						for i, key := range keys {
							id := &Identifier{scope: id.scope, key: key}
							insert(values, id, strings.Replace(xs[i], ",\\ ", ", ", -1))
						}
					}
				}
			}
		}
	}

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
