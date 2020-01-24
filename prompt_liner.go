package main

import (
	"io"

	"github.com/itchyny/liner"
	"github.com/mattn/go-tty"
)

type realPrompt struct {
	state *liner.State
}

func newPrompt() *realPrompt {
	return &realPrompt{}
}

func (p *realPrompt) start() {
	tty, err := tty.Open()
	if err != nil {
		panic(err)
	}
	p.state = liner.NewLinerTTY(tty)
	p.state.SetCtrlCAborts(true)
}

func (p *realPrompt) prompt(message string) (string, error) {
	input, err := p.state.Prompt(message)
	if err != nil {
		if err == liner.ErrPromptAborted {
			return "", io.EOF
		}
		return "", err
	}
	return input, nil
}

func (p *realPrompt) setHistory(history []string) {
	p.state.ClearHistory()
	for i := len(history) - 1; i >= 0; i-- {
		p.state.AppendHistory(history[i])
	}
}

func (p *realPrompt) close() {
	p.state.Close()
}
