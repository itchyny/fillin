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
	setHistory := func(history []string) {
		line.ClearHistory()
		for i := len(history) - 1; i >= 0; i-- {
			line.AppendHistory(history[i])
		}
	}
	checkErr := func(err error) {
		if err != nil {
			if err == liner.ErrPromptAborted || err == io.EOF {
				os.Exit(1)
			}
			log.Fatal(err)
		}
	}
	values := make(map[string]map[string]string)

	scopeAsked := make(map[string]bool)
	if in == nil {
		for _, id := range identifiers {
			if found(values, id) || id.scope == "" || scopeAsked[id.scope] {
				continue
			}
			scopeAsked[id.scope] = true
			idg := collect(identifiers, id.scope)
			if len(idg.keys) == 0 {
				continue
			}
			history := config.collectScopedPairHistory(idg)
			if len(history) == 0 {
				continue
			}
			setHistory(history)
			prompt := fmt.Sprintf("[%s] %s: ", id.scope, strings.Join(idg.keys, ", "))
			text, err := line.Prompt(prompt)
			checkErr(err)
			xs := strings.Split(strings.TrimSuffix(text, "\n"), ", ")
			if len(xs) == len(idg.keys) {
				for i, key := range idg.keys {
					insert(values, &Identifier{scope: id.scope, key: key}, strings.Replace(xs[i], ",\\ ", ", ", -1))
				}
			}
		}
	}

	for _, id := range identifiers {
		if found(values, id) {
			continue
		}
		prompt := fmt.Sprintf("%s: ", id.key)
		if id.scope != "" {
			prompt = fmt.Sprintf("[%s] %s: ", id.scope, id.key)
		}
		var text string
		var err error
		if in == nil {
			setHistory(config.collectHistory(id))
			text, err = line.Prompt(prompt)
			checkErr(err)
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

	return values
}
