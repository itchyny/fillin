package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/peterh/liner"
)

// Resolve asks the user to resolve the identifiers
func Resolve(identifiers []*Identifier, config *Config, in io.Reader, out io.Writer) (map[string]map[string]string, error) {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	setHistory := func(history []string) {
		line.ClearHistory()
		for i := len(history) - 1; i >= 0; i-- {
			line.AppendHistory(history[i])
		}
	}
	normalizeErr := func(err error) error {
		if err == liner.ErrPromptAborted || err == io.EOF {
			return io.EOF
		}
		return err
	}
	values := make(map[string]map[string]string)
	var bin *bufio.Reader
	if in != nil {
		bin = bufio.NewReader(in)
	}

	scopeAsked := make(map[string]bool)
	for _, id := range identifiers {
		if found(values, id) || id.scope == "" || scopeAsked[id.scope] {
			continue
		}
		scopeAsked[id.scope] = true
		idg := collect(identifiers, id.scope)
		if len(idg.keys) == 0 {
			continue
		}
		var text string
		var err error
		history := config.collectScopedPairHistory(idg)
		if len(history) == 0 {
			continue
		}
		if bin == nil {
			setHistory(history)
			text, err = line.Prompt(idg.prompt())
			if err := normalizeErr(err); err != nil {
				return nil, err
			}
		} else {
			out.Write([]byte(idg.prompt()))
			text, err = bin.ReadString('\n')
			if err != nil {
				return nil, err
			}
		}
		xs := strings.Split(strings.TrimSuffix(text, "\n"), ", ")
		if len(xs) == len(idg.keys) {
			for i, key := range idg.keys {
				insert(values, &Identifier{scope: id.scope, key: key}, strings.Replace(xs[i], ",\\ ", ", ", -1))
			}
		}
	}

	for _, id := range identifiers {
		if found(values, id) {
			continue
		}
		var text string
		var err error
		if bin == nil {
			setHistory(config.collectHistory(id))
			text, err = line.Prompt(id.prompt())
			if err := normalizeErr(err); err != nil {
				return nil, err
			}
		} else {
			out.Write([]byte(id.prompt()))
			text, err = bin.ReadString('\n')
			if err != nil {
				return nil, err
			}
		}
		insert(values, id, strings.TrimSuffix(text, "\n"))
	}

	return values, nil
}
