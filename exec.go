package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	cmd, err := Run(configDir, os.Args[1:], newPrompt())
	if err != nil {
		return err
	}
	if err := syscallExec(sh, append(cmdBase, cmd), os.Environ()); err != nil {
		return err
	}
	return nil
}

func getConfigDir() (string, error) {
	if dir := os.Getenv("FILLIN_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, name), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", name), nil
}

func printVersion() {
	fmt.Printf("%s %s\n", name, version)
}

func printHelp() {
	fmt.Printf(`NAME:
   %[1]s - %[2]s

USAGE:
   %[1]s command...

EXAMPLES:
   %[1]s echo {{message}} # in bash/zsh shell
   %[1]s echo [[message]] # in fish shell
   %[1]s psql -h {{psql:hostname}} -U {{psql:username}} -d {{psql:dbname}}
   %[1]s curl {{example-api:base-url}}/api/1/example/info -H 'Authorization: Bearer {{example-api:access-token}}'

VERSION:
   %[3]s

AUTHOR:
   %[4]s
`, name, description, version, author)
}
