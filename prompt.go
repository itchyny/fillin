package main

type prompt interface {
	start() error
	prompt(string) (string, error)
	setHistory([]string)
	close()
}
