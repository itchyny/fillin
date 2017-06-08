package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Fillin(args []string) []string {
	ret := make([]string, len(args))
	reader := bufio.NewReader(os.Stdin)
	for i, arg := range args {
		ret[i] = fillinStr(arg, reader)
	}
	return ret
}

var fillinPattern = regexp.MustCompile(`{{[-0-9A-Za-z_]+}}`)

func fillinStr(str string, reader *bufio.Reader) string {
	return fillinPattern.ReplaceAllStringFunc(str, func(match string) string {
		identifier := match[2 : len(match)-2]
		fmt.Printf("%s: ", identifier)
		text, _ := reader.ReadString('\n')
		return strings.TrimSuffix(text, "\n")
	})
}
