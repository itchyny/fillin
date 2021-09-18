package main

import (
	"encoding/json"
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
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "fillin.json")
	config, err := readConfig(path)
	if err != nil {
		return "", err
	}
	filled, err := Fillin(args, config, p)
	if err != nil {
		return "", err
	}
	cmd := escapeJoin(filled)
	if err := writeConfig(path, config); err != nil {
		return "", err
	}
	if err := appendHistory(dir, cmd); err != nil {
		return "", err
	}
	return cmd, nil
}

func readConfig(path string) (*Config, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0o600)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func writeConfig(path string, config *Config) error {
	tmp, err := ioutil.TempFile(filepath.Dir(path), "fillin-*.json")
	if err != nil {
		return err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	tmp.Close()
	return os.Rename(tmp.Name(), path)
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

func appendHistory(dir, cmd string) error {
	if cmd == "" {
		return nil
	}
	histfile := filepath.Join(dir, ".fillin.histfile")
	f, err := os.OpenFile(histfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		f.Chmod(0o600)
		f.Close()
	}()
	zshhist.NewWriter(f).Write(
		zshhist.History{Time: int(time.Now().Unix()), Elapsed: 0, Command: cmd},
	)
	return f.Sync()
}
