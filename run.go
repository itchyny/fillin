package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/zshhist-go"
)

// Run fillin
func Run(dir string, args []string, p prompt) (string, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return "", err
	}

	path := filepath.Join(dir, "fillin.json")
	rfile, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	defer rfile.Close()
	w := new(bytes.Buffer)
	filled, err := Fillin(args, rfile, w, p)
	if err != nil {
		return "", err
	}
	cmd := escapeJoin(filled)
	rfile.Close() // not be defered due to rename

	tmp, err := ioutil.TempFile(dir, "fillin-*.json")
	if err != nil {
		return "", err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	if n, err := tmp.Write(w.Bytes()); n != w.Len() || err != nil {
		return "", err
	}
	if err := tmp.Sync(); err != nil {
		return "", err
	}
	tmp.Close() // not be defered due to rename

	if err := os.Rename(tmp.Name(), path); err != nil {
		return "", err
	}

	if cmd != "" {
		histfile := filepath.Join(dir, ".fillin.histfile")
		hfile, err := os.OpenFile(histfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return "", err
		}
		defer hfile.Close()
		w := zshhist.NewWriter(hfile)
		w.Write(zshhist.History{Time: int(time.Now().Unix()), Elapsed: 0, Command: cmd})
		hfile.Chmod(0600)
		if err := hfile.Sync(); err != nil {
			return "", err
		}
	}

	return cmd, nil
}

func escapeJoin(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	for i, arg := range args {
		args[i] = escape(arg)
	}
	return strings.Join(args, " ")
}

var redirectPattern = regexp.MustCompile(`^\s*[012]?\s*[<>]`)

func escape(arg string) string {
	switch arg {
	case "|", "||", "&&", ">", ">>", "<":
		return arg
	}
	if redirectPattern.MatchString(arg) {
		return arg
	}
	if !strings.ContainsAny(arg, "|&><[?! \"'\a\b\f\n\r\t\v\\") {
		return arg
	}
	return strconv.Quote(arg)
}
