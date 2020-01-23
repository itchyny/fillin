package main

type prompt interface {
	start()
	prompt(string) (string, error)
	setHistory([]string)
	close()
}
