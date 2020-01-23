package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type testPrompt struct {
	reader *bufio.Reader
	writer io.Writer
}

func newTestPrompt(input string) *testPrompt {
	return &testPrompt{
		reader: bufio.NewReader(strings.NewReader(input)),
		writer: new(bytes.Buffer),
	}
}

func (p *testPrompt) start() {
}

func (p *testPrompt) prompt(message string) (string, error) {
	p.writer.Write([]byte(message))
	return p.reader.ReadString('\n')
}

func (p *testPrompt) setHistory([]string) {
}

func (p *testPrompt) close() {
}
