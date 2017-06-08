package main

import (
	"regexp"
)

var fillinPattern = regexp.MustCompile(`{{[-0-9A-Za-z_]+}}`)

func Fillin(args []string) []string {
	ret := make([]string, len(args))
	var identifiers []string
	for _, arg := range args {
		matches := fillinPattern.FindAllString(arg, -1)
		for _, match := range matches {
			identifiers = append(identifiers, identifierFromMatch(match))
		}
	}
	values := Resolve(identifiers)
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
