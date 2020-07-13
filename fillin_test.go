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

var mapCompareTests = []struct {
	m1, m2   map[string]string
	expected int
}{
	{
		m1:       map[string]string{},
		m2:       map[string]string{},
		expected: mapCompareEqual,
	},
	{
		m1:       map[string]string{"x": "1"},
		m2:       map[string]string{},
		expected: mapCompareSuperset,
	},
	{
		m1:       map[string]string{"x": "1", "y": "2"},
		m2:       map[string]string{},
		expected: mapCompareSuperset,
	},
	{
		m1:       map[string]string{},
		m2:       map[string]string{"x": "1"},
		expected: mapCompareSubset,
	},
	{
		m1:       map[string]string{"x": "1"},
		m2:       map[string]string{"x": "1", "y": "2"},
		expected: mapCompareSubset,
	},
	{
		m1:       map[string]string{"x": "1", "y": "2"},
		m2:       map[string]string{"x": "1"},
		expected: mapCompareSuperset,
	},
	{
		m1:       map[string]string{"x": "1"},
		m2:       map[string]string{"y": "2"},
		expected: mapCompareDiff,
	},
	{
		m1:       map[string]string{"x": "1"},
		m2:       map[string]string{"x": "2"},
		expected: mapCompareDiff,
	},
	{
		m1:       map[string]string{"x": "1", "y": "2", "z": "3"},
		m2:       map[string]string{"x": "1", "w": "2", "z": "3"},
		expected: mapCompareDiff,
	},
	{
		m1:       map[string]string{"x": "1"},
		m2:       map[string]string{"x": "1", "y": "2", "z": "3"},
		expected: mapCompareSubset,
	},
	{
		m1:       map[string]string{"x": "1", "y": "2", "z": "3"},
		m2:       map[string]string{"x": "1", "y": "2", "z": "3"},
		expected: mapCompareEqual,
	},
}

func Test_mapCompare(t *testing.T) {
	for _, test := range mapCompareTests {
		got := mapCompare(test.m1, test.m2)
		if got != test.expected {
			t.Errorf("mapCompare failed for %+v", test)
		}
	}
}
