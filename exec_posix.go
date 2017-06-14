// +build !windows

package main

import "syscall"

func syscallExec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}
