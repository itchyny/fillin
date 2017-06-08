package main

import (
	"os/exec"
	"syscall"

	"github.com/kballard/go-shellquote"
)

func run(args []string, env []string) error {
	sh, err := exec.LookPath("sh")
	if err != nil {
		return err
	}
	cmd := []string{"sh", "-c", shellquote.Join(args...)}
	if err := syscall.Exec(sh, cmd, env); err != nil {
		return err
	}
	return nil
}
