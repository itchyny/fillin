package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

var fillinPattern = regexp.MustCompile(`{{\s*[-0-9A-Za-z_]+(\s*:\s*[-0-9A-Za-z_]+)?\s*}}`)

func collectIdentifiers(args []string) []*Identifier {
	var identifiers []*Identifier
	for _, arg := range args {
		matches := fillinPattern.FindAllString(arg, -1)
		for _, match := range matches {
			identifiers = append(identifiers, identifierFromMatch(match))
		}
	}
	return identifiers
}

// Fillin fills in the command arguments
func Fillin(args []string, r io.Reader, w io.Writer, in *bufio.Reader, out *bufio.Writer) []string {
	ret := make([]string, len(args))
	config, err := ReadConfig(r)
	if err != nil {
		log.Fatal(err)
	}
	if config.Scopes == nil {
		config.Scopes = make(map[string]*Scope)
	}
	values := Resolve(collectIdentifiers(args), config, in, out)
	insertValues(config.Scopes, values)
	if err := WriteConfig(w, config); err != nil {
		log.Fatal(err)
	}
	for i, arg := range args {
		ret[i] = fillinPattern.ReplaceAllStringFunc(arg, func(match string) string {
			return lookup(values, identifierFromMatch(match))
		})
	}
	return ret
}

func identifierFromMatch(match string) *Identifier {
	match = match[2 : len(match)-2]
	var scope, key string
	if strings.ContainsRune(match, ':') {
		xs := strings.Split(match, ":")
		scope = strings.TrimSpace(xs[0])
		key = strings.TrimSpace(xs[1])
	} else {
		key = strings.TrimSpace(match)
	}
	return &Identifier{scope, key}
}

func insertValues(scopes map[string]*Scope, values map[string]map[string]string) {
	for scope := range values {
		if _, ok := scopes[scope]; !ok {
			scopes[scope] = &Scope{}
		}
		newValues := make([]map[string]string, 0)
		strs := make(map[string]bool)
		insert := func(v map[string]string) {
			s := stringifyValue(v)
			if _, ok := strs[s]; !ok {
				strs[s] = true
				newValues = append(newValues, v)
			}
		}
		insert(values[scope])
		for _, v := range scopes[scope].Values {
			insert(v)
		}
		scopes[scope].Values = newValues
	}
}
