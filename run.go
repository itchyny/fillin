package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/itchyny/zshhist-go"
)

// Run fillin
func Run(dir string, args []string, in io.Reader, out io.Writer) (string, error) {
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
	w := new(bytes.Buffer)
	filled, err := Fillin(args, rfile, w, in, out)
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
	tmp.Close() // not be defered due to rename

	if err := os.Rename(tmp.Name(), path); err != nil {
		return "", err
	}

	if cmd != "" {
		histfile := filepath.Join(dir, ".fillin.histfile")
		hfile, err := os.OpenFile(histfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		defer hfile.Close()
		if err != nil {
			return "", err
		}
		w := zshhist.NewWriter(hfile)
		w.Write(zshhist.History{Time: int(time.Now().Unix()), Elapsed: 0, Command: cmd})
		hfile.Chmod(0600)
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
