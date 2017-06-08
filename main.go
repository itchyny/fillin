package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/kballard/go-shellquote"
)

func main() {
	sh, err := exec.LookPath("sh")
	if err != nil {
		log.Fatal(err)
	}
	cmd := []string{"sh", "-c", shellquote.Join(os.Args[1:]...)}
	env := os.Environ()
	if err := syscall.Exec(sh, cmd, env); err != nil {
		log.Fatal(err)
	}
}
