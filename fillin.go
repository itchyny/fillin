package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
)

var fillinPattern = regexp.MustCompile(`{{[-0-9A-Za-z_]+}}`)

// Fillin fills in the command arguments
func Fillin(args []string, r io.Reader, w io.Writer, in *bufio.Reader) []string {
	ret := make([]string, len(args))
	var identifiers []string
	for _, arg := range args {
		matches := fillinPattern.FindAllString(arg, -1)
		for _, match := range matches {
			identifiers = append(identifiers, identifierFromMatch(match))
		}
	}
	config, err := ReadConfig(r)
	if err != nil {
		log.Fatal(err)
	}
	if config.Scopes == nil {
		config.Scopes = make(map[string]*Scope)
	}
	values := Resolve(identifiers, config, in)
	scope := ""
	if _, ok := config.Scopes[scope]; !ok {
		config.Scopes[scope] = &Scope{}
	}
	config.Scopes[scope].Values = insertValue(config.Scopes[scope].Values, values)
	if err := WriteConfig(w, config); err != nil {
		log.Fatal(err)
	}
	for i, arg := range args {
		ret[i] = fillinPattern.ReplaceAllStringFunc(arg, func(match string) string {
			return values[identifierFromMatch(match)]
		})
	}
	return ret
}

func identifierFromMatch(match string) string {
	return match[2 : len(match)-2]
}

func insertValue(orig []map[string]string, value map[string]string) []map[string]string {
	values := make([]map[string]string, 0, len(orig)+1)
	strs := make(map[string]bool)
	insert := func(v map[string]string) {
		s := stringifyValue(v)
		if _, ok := strs[s]; !ok {
			strs[s] = true
			values = append(values, v)
		}
	}
	insert(value)
	for _, v := range orig {
		insert(v)
	}
	return values
}
