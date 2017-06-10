package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kballard/go-shellquote"
)

var homedir string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homedir = usr.HomeDir
}

// Run fillin
func Run(configPath string, args []string, in *bufio.Reader) (string, error) {
	path := filepath.Join(strings.Split(configPath, "/")...)
	if path[0] == '~' {
		path = homedir + path[1:]
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	rfile, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	w := new(bytes.Buffer)
	cmd := shellquote.Join(Fillin(args, rfile, w, in)...)
	rfile.Close() // not be defered due to rename
	tmpFileName := fmt.Sprintf("fillin.%d-%d.json", os.Getpid(), rand.Int())
	tmp := filepath.Join(filepath.Dir(path), tmpFileName)
	defer os.Remove(tmp)
	wfile, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	if n, err := wfile.Write(w.Bytes()); n != w.Len() || err != nil {
		wfile.Close()
		return "", err
	}
	wfile.Close() // not be defered due to rename
	if err := os.Rename(tmp, path); err != nil {
		return "", err
	}
	return cmd, nil
}
