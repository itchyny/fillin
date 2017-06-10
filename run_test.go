package main

import (
	"bufio"
	"bytes"
	"reflect"
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
		args: []string{"echo", "{{foo}},", "{{bar}},", "{{foo}}-{{bar}}-{{baz}}"},
		in: `Hello
wonderful
world!
`,
		expected: `echo Hello, wonderful, Hello-wonderful-world\!`,
	},
}

func TestRun(t *testing.T) {
	path := "./.test/run.json"
	for _, test := range runTests {
		in := bufio.NewReader(bytes.NewBufferString(test.in))
		cmd, err := Run(path, test.args, in)
		if err != nil {
			t.Errorf("error occurred unexpectedly: %+v", err)
		}
		if !reflect.DeepEqual(cmd, test.expected) {
			t.Errorf("command not correct (expected: %+v, got: %+v)", test.expected, cmd)
		}
	}
}
