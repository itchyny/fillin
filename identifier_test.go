package main

import (
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

func Test_found(t *testing.T) {
	for _, test := range identifierTests {
		got := found(test.values, &test.id)
		if got != test.found {
			t.Errorf("found not correct for %+v (found: %+v, got: %+v)", test.id, test.found, got)
		}
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

func Test_lookup(t *testing.T) {
	for _, test := range identifierTests {
		got := lookup(test.values, &test.id)
		if got != test.value {
			t.Errorf("lookup not correct for %+v (found: %+v, got: %+v)", test.id, test.value, got)
		}
	}
}
