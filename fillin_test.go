package main

import (
	"reflect"
	"testing"
)

var collectIdentifiersTests = []struct {
	args []string
	ids  []*Identifier
}{
	{
		args: []string{"{{foo}}", "{{bar}}", "{{baz}}"},
		ids: []*Identifier{
			{key: "foo"},
			{key: "bar"},
			{key: "baz"},
		},
	},
	{
		args: []string{"{{foo:bar}}", "{{foo:baz}}", "{{foo}}"},
		ids: []*Identifier{
			{scope: "foo", key: "bar"},
			{scope: "foo", key: "baz"},
			{scope: "", key: "foo"},
		},
	},
	{
		args: []string{"[[foo]]", "{{bar}}", "[[baz]]"},
		ids: []*Identifier{
			{key: "foo"},
			{key: "bar"},
			{key: "baz"},
		},
	},
	{
		args: []string{"[[foo:bar]]", "[[foo:baz]]", "{{foo}}"},
		ids: []*Identifier{
			{scope: "foo", key: "bar"},
			{scope: "foo", key: "baz"},
			{scope: "", key: "foo"},
		},
	},
}

func Test_collectIdentifiers(t *testing.T) {
	for _, test := range collectIdentifiersTests {
		got := collectIdentifiers(test.args)
		if !reflect.DeepEqual(got, test.ids) {
			t.Errorf("collectIdentifiers failed for %+v", test.args)
		}
	}
}
