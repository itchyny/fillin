// +build windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func syscallExec(argv0 string, argv []string, envv []string) error {
	cmd := exec.Command(argv0, argv[1:]...)
	cmd.Env = envv
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil && cmd.ProcessState == nil {
		return err
	}
	os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
	panic("unreachable")
}
