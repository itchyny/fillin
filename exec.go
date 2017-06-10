package main

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

var cmdBase = []string{"sh", "-c"}

func init() {
	if runtime.GOOS == "windows" {
		cmdBase = []string{"cmd", "/c"}
	}
}

func Exec() error {
	sh, err := exec.LookPath(cmdBase[0])
	if err != nil {
		return err
	}
	configPath := "~/.config/fillin/fillin.json"
	in := bufio.NewReader(os.Stdin)
	cmd, err := Run(configPath, os.Args[1:], in)
	if err != nil {
		return err
	}
	if err := syscall.Exec(sh, append(cmdBase, cmd), os.Environ()); err != nil {
		return err
	}
	return nil
}
