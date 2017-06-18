package main

import (
	"reflect"
	"testing"
)

var identifierTests = []struct {
	values map[string]map[string]string
	id     Identifier
	found  bool
	value  string
}{
	{
		values: nil,
		id:     Identifier{key: "foo"},
		found:  false,
	},
	{
		values: map[string]map[string]string{"": {"foo": "example"}},
		id:     Identifier{key: "foo"},
		found:  true,
		value:  "example",
	},
	{
		values: map[string]map[string]string{"": {"foo": "example"}},
		id:     Identifier{key: "bar"},
		found:  false,
	},
	{
		values: map[string]map[string]string{"example": {"foo": "example"}},
		id:     Identifier{scope: "example", key: "foo"},
		found:  true,
		value:  "example",
	},
	{
		values: map[string]map[string]string{"": {"foo": "example"}},
		id:     Identifier{scope: "example", key: "foo"},
		found:  false,
	},
	{
		values: map[string]map[string]string{"example": {"foo": "example"}},
		id:     Identifier{scope: "example", key: "bar"},
		found:  false,
	},
}

func Test_prompt(t *testing.T) {
	id1 := &Identifier{key: "foo"}
	got := id1.prompt()
	if got != "foo: " {
		t.Errorf("prompt is not correct for %+v (found: %+v, got: %+v)", id1, "foo: ", got)
	}
	id2 := &Identifier{scope: "foo", key: "bar"}
	got = id2.prompt()
	if got != "[foo] bar: " {
		t.Errorf("prompt is not correct for %+v (found: %+v, got: %+v)", id2, "[foo] bar: ", got)
	}
	idg := &IdentifierGroup{scope: "foo", keys: []string{"bar", "baz", "qux"}}
	got = idg.prompt()
	if got != "[foo] bar, baz, qux: " {
		t.Errorf("prompt is not correct for %+v (found: %+v, got: %+v)", idg, "[foo] bar, baz, qux: ", got)
	}
}

func Test_found(t *testing.T) {
	for _, test := range identifierTests {
		got := found(test.values, &test.id)
		if got != test.found {
			t.Errorf("found not correct for %+v (found: %+v, got: %+v)", test.id, test.found, got)
		}
	}
}

func Test_collect(t *testing.T) {
	ids := []*Identifier{
		{scope: "foo", key: "foo"},
		{scope: "foo", key: "bar"},
		{scope: "zoo", key: "foo"},
		{scope: "foo", key: "foo"},
		{scope: "foo", key: "baz"},
		{scope: "qux", key: "bar"},
	}
	expectedFoo := &IdentifierGroup{
		scope: "foo",
		keys:  []string{"foo", "bar", "baz"},
	}
	expectedBar := &IdentifierGroup{
		scope: "bar",
		keys:  nil,
	}
	idgFoo := collect(ids, "foo")
	if !reflect.DeepEqual(idgFoo, expectedFoo) {
		t.Errorf("collect not correct (expected: %+v, got: %+v)", expectedFoo, idgFoo)
	}
	idgBar := collect(ids, "bar")
	if !reflect.DeepEqual(idgBar, expectedBar) {
		t.Errorf("collect not correct (expected: %+v, got: %+v)", expectedBar, idgBar)
	}
}

func Test_insert(t *testing.T) {
	values := make(map[string]map[string]string)
	id := &Identifier{key: "foo"}
	value := "bar"
	insert(values, id, value)
	v, ok := values[""]["foo"]
	if !ok {
		t.Errorf("insert failed for %+v", id)
	}
	if v != value {
		t.Errorf("insert not correctly for %+v (found: %+v, got: %+v)", id, v, value)
	}
	id = &Identifier{scope: "foo", key: "bar"}
	value = "example"
	insert(values, id, value)
	v, ok = values["foo"]["bar"]
	if !ok {
		t.Errorf("insert failed for %+v", id)
	}
	if v != value {
		t.Errorf("insert not correctly for %+v (found: %+v, got: %+v)", id, v, value)
	}
}

func Test_empty(t *testing.T) {
	tests := []struct {
		values   map[string]map[string]string
		expected bool
	}{
		{
			values:   nil,
			expected: true,
		},
		{
			values: map[string]map[string]string{
				"foo": {},
			},
			expected: true,
		},
		{
			values: map[string]map[string]string{
				"foo": {
					"bar": "",
					"baz": "",
				},
			},
			expected: true,
		},
		{
			values: map[string]map[string]string{
				"foo": {
					"bar": "",
					"baz": "",
				},
				"bar": {
					"qux": "quux",
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		got := empty(test.values)
		if got != test.expected {
			t.Errorf("empty not correctly for %+v (expected: %+v, got: %+v)", test.values, test.expected, got)
		}
	}
}

func Test_lookup(t *testing.T) {
	for _, test := range identifierTests {
		got := lookup(test.values, &test.id)
		if got != test.value {
			t.Errorf("lookup not correct for %+v (expected: %+v, got: %+v)", test.id, test.value, got)
		}
	}
}
