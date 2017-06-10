package main

import (
	"bufio"
	"bytes"
	"reflect"
	"sync"
	"testing"
)

var runTests = []struct {
	args     []string
	in       string
	expected string
}{
	{
		args:     []string{"echo", "Hello,", "world!"},
		in:       ``,
		expected: `echo Hello, world\!`,
	},
	{
		args: []string{"echo", "{{foo}},", "{{bar}}"},
		in: `Hello test
world test!
`,
		expected: `echo 'Hello test,' 'world test!'`,
	},
	{
		args: []string{"echo", "{{ foo-bar_baz }}", "{{ FOO-9 : BAR_2 }}", "{{ X }}"},
		in: `Foo bar
FOO BAR
X
`,
		expected: `echo 'Foo bar' 'FOO BAR' X`,
	},
	{
		args: []string{"echo", "{{foo}},", "{{bar}},", "{{foo}}-{{bar}}-{{baz}}"},
		in: `Hello
wonderful
world!
`,
		expected: `echo Hello, wonderful, Hello-wonderful-world\!`,
	},
	{
		args: []string{"echo", "{{foo:bar}},", "{{ foo : bar }},", "{{foo:baz}}"},
		in: `Hello
example world!
`,
		expected: `echo Hello, Hello, 'example world!'`,
	},
	{
		args: []string{"echo", "{{foo}},", "{{bar}}", "{{baz}}"},
		in: `こんにちは
世界
+。:.ﾟ٩(๑>◡<๑)۶:.｡+ﾟ
`,
		expected: `echo こんにちは, 世界 +。:.ﾟ٩\(๑\>◡\<๑\)۶:.｡+ﾟ`,
	},
	{
		args: []string{"echo", "{{foo}}", "|", "echo", "||", "echo", "&&", "echo", ">", "/dev/null", "</dev/null", "2>&1", "1", ">&2"},
		in: `Hello world!
`,
		expected: `echo 'Hello world!' | echo || echo && echo > /dev/null </dev/null 2>&1 1 >&2`,
	},
	{
		args: []string{"echo", "{{foo}}", "{{bar}}"},
		in: `\'"${[]}|&;<>()*?!
	foo bar baz
`,
		expected: `echo \\\'\"\$\{\[]}\|\&\;\<\>\(\)\*\?\! '\tfoo bar baz'`,
	},
}

func TestRun(t *testing.T) {
	path := "./.test/run.json"
	for _, test := range runTests {
		in := bufio.NewReader(bytes.NewBufferString(test.in))
		out := bufio.NewWriter(new(bytes.Buffer))
		cmd, err := Run(path, test.args, in, out)
		if err != nil {
			t.Errorf("error occurred unexpectedly: %+v", err)
		}
		if !reflect.DeepEqual(cmd, test.expected) {
			t.Errorf("command not correct (expected: %+v, got: %+v)", test.expected, cmd)
		}
	}
}

func TestRun_concurrently(t *testing.T) {
	path := "./.test/concurrently.json"
	test := runTests[1]
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			in := bufio.NewReader(bytes.NewBufferString(test.in))
			out := bufio.NewWriter(new(bytes.Buffer))
			cmd, err := Run(path, test.args, in, out)
			if err != nil {
				t.Errorf("error occurred unexpectedly: %+v", err)
			}
			if !reflect.DeepEqual(cmd, test.expected) {
				t.Errorf("command not correct (expected: %+v, got: %+v)", test.expected, cmd)
			}
		}()
	}
	wg.Wait()
}
