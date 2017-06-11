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
	"time"
	"unicode/utf8"
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
func Run(configPath string, args []string, in *bufio.Reader, out *bufio.Writer) (string, error) {
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
	cmd := escapeJoin(Fillin(args, rfile, w, in, out))
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
	histfile := filepath.Join(filepath.Dir(path), ".fillin.histfile")
	hfile, err := os.OpenFile(histfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer hfile.Close()
	if err != nil {
		return "", err
	}
	hfile.WriteString(fmt.Sprintf(": %d:0;%s\n", time.Now().Unix(), cmd))
	return cmd, nil
}

func escapeJoin(args []string) string {
	for i, arg := range args {
		args[i] = escape(arg)
	}
	return strings.Join(args, " ")
}

func escape(arg string) string {
	switch arg {
	case "|", "||", "&&", ">", ">>", "<":
		return arg
	}
quote:
	for _, quote := range []bool{false, true} {
		s := arg
		isHead, afterFd, afterRedirect := true, false, false
		var buf bytes.Buffer
		if quote {
			buf.WriteByte('\'')
		}
		for len(s) > 0 {
			c, l := utf8.DecodeRuneInString(s)
			s = s[l:]
			if (isHead || afterFd) && strings.ContainsRune("<>", c) && !quote {
				buf.WriteRune(c)
				isHead, afterFd, afterRedirect = false, false, true
				continue
			} else if afterRedirect && strings.ContainsRune("<>", c) && !quote {
				buf.WriteRune(c)
			} else if isHead && strings.ContainsRune("12", c) && !quote {
				buf.WriteRune(c)
				isHead, afterFd = false, true
				continue
			} else if afterRedirect && strings.ContainsRune("&", c) && !quote {
				buf.WriteRune(c)
			} else if strings.ContainsRune("\\'\"`${[|&;<>()*?!", c) && !quote {
				buf.WriteByte('\\')
				buf.WriteRune(c)
			} else if c == rune(' ') && !quote {
				continue quote
			} else if c == rune('\t') {
				if quote {
					buf.WriteByte('\\')
					buf.WriteByte('t')
				} else {
					continue quote
				}
			} else {
				buf.WriteRune(c)
			}
			isHead, afterFd, afterRedirect = false, false, false
		}
		if quote {
			buf.WriteByte('\'')
		}
		return buf.String()
	}
	return ""
}
