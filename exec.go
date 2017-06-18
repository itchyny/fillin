package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var cmdBase = []string{"sh", "-c"}

func init() {
	if runtime.GOOS == "windows" {
		cmdBase = []string{"cmd", "/c"}
	}
}

// Exec fillin
func Exec() error {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "-v", "version", "-version", "--version":
			printVersion()
			return nil
		case "-h", "help", "-help", "--help":
			printHelp()
			return nil
		}
	}
	sh, err := exec.LookPath(cmdBase[0])
	if err != nil {
		return err
	}
	configDir := "~/.config/fillin"
	cmd, err := Run(configDir, os.Args[1:], nil, bufio.NewWriter(os.Stdout))
	if err != nil {
		return err
	}
	if err := syscallExec(sh, append(cmdBase, cmd), os.Environ()); err != nil {
		return err
	}
	return nil
}

func printVersion() {
	fmt.Printf("%s version %s\n", name, version)
}

func printHelp() {
	fmt.Printf(strings.Replace(`NAME:
   $NAME - %s

USAGE:
   $NAME command...

EXAMPLES:
   $NAME echo {{message}} # in bash/zsh shell
   $NAME echo [[message]] # in fish shell
   $NAME psql -h {{psql:hostname}} -U {{psql:username}} -d {{psql:dbname}}
   $NAME curl {{example-api:base-url}}/api/1/example/info -H 'Authorization: Bearer {{example-api:access-token}}'

VERSION:
   %s

AUTHOR:
   %s
`, "$NAME", name, -1), description, version, author)
}
