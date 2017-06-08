package main

import (
	"os/exec"
	"runtime"
	"syscall"

	"github.com/kballard/go-shellquote"
)

var cmdBase = []string{"sh", "-c"}

func init() {
	if runtime.GOOS == "windows" {
		cmdBase = []string{"cmd", "/c"}
	}
}

func run(args []string, env []string) error {
	sh, err := exec.LookPath(cmdBase[0])
	if err != nil {
		return err
	}
	cmd := append(cmdBase, shellquote.Join(Fillin(args)...))
	if err := syscall.Exec(sh, cmd, env); err != nil {
		return err
	}
	return nil
}
