package main

import (
	"bytes"
	"reflect"
	"testing"
)

var configTests = []struct {
	config Config
	bytes  []byte
}{
	{
		config: Config{
			Scopes: map[string]*Scope{},
		},
		bytes: []byte(`{
  "scopes": {}
}
`),
	},
	{
		config: Config{
			Scopes: map[string]*Scope{
				"": &Scope{},
			},
		},
		bytes: []byte(`{
  "scopes": {
    "": {
      "values": null
    }
  }
}
`),
	},
	{
		config: Config{
			Scopes: map[string]*Scope{
				"": &Scope{
					Values: []map[string]string{
						{
							"baz": "world!",
							"foo": "Hello 1",
						},
						{
							"baz": "world!",
							"foo": "Hello 2",
						},
					},
				},
				"sample": &Scope{
					Values: []map[string]string{
						{
							"foo": "Test1, world!",
							"bar": "test1, test",
							"baz": "baz",
						},
						{
							"foo": "Test2, world!",
							"bar": "test2, test",
						},
					},
				},
			},
		},
		bytes: []byte(`{
  "scopes": {
    "": {
      "values": [
        {
          "baz": "world!",
          "foo": "Hello 1"
        },
        {
          "baz": "world!",
          "foo": "Hello 2"
        }
      ]
    },
    "sample": {
      "values": [
        {
          "bar": "test1, test",
          "baz": "baz",
          "foo": "Test1, world!"
        },
        {
          "bar": "test2, test",
          "foo": "Test2, world!"
        }
      ]
    }
  }
}
`),
	},
}

func TestReadConfig(t *testing.T) {
	for _, test := range configTests {
		r := bytes.NewReader(test.bytes)
		config, err := ReadConfig(r)
		if err != nil {
			t.Errorf("error occurred unexpectedly on reading a config %+v", err)
		}
		if !reflect.DeepEqual(test.config, *config) {
			t.Errorf("config loaded incorrectly (expected: %+v, got: %+v)", test.config, *config)
		}
	}
}

func TestWriteConfig(t *testing.T) {
	for _, test := range configTests {
		w := new(bytes.Buffer)
		err := WriteConfig(w, &test.config)
		if err != nil {
			t.Errorf("error occurred unexpectedly on writing a config %+v", err)
		}
		if !reflect.DeepEqual(test.bytes, w.Bytes()) {
			t.Errorf("config wrote incorrectly (expected: %+v, got: %+v)", string(test.bytes), w.String())
		}
	}
}

func Test_collectHistory(t *testing.T) {
	hs1 := configTests[2].config.collectHistory(&Identifier{key: "foo"})
	expected1 := []string{"Hello 1", "Hello 2"}
	if !reflect.DeepEqual(hs1, expected1) {
		t.Errorf("collectHistory incorrect (expected: %+v, got: %+v)", expected1, hs1)
	}
	hs2 := configTests[2].config.collectHistory(&Identifier{scope: "sample", key: "foo"})
	expected2 := []string{"Test1, world!", "Test2, world!"}
	if !reflect.DeepEqual(hs2, expected2) {
		t.Errorf("collectHistory incorrect (expected: %+v, got: %+v)", expected2, hs2)
	}
	hs3 := configTests[2].config.collectHistory(&Identifier{scope: "foo", key: "test"})
	expected3 := []string{}
	if !reflect.DeepEqual(hs3, expected3) {
		t.Errorf("collectHistory incorrect (expected: %+v, got: %+v)", expected3, hs3)
	}
}

func Test_collectScopedPairHistory(t *testing.T) {
	hs1 := configTests[2].config.collectScopedPairHistory(&IdentifierGroup{keys: []string{"foo", "baz"}})
	expected1 := []string{"Hello 1, world!", "Hello 2, world!"}
	if !reflect.DeepEqual(hs1, expected1) {
		t.Errorf("collectScopedPairHistory incorrect (expected: %+v, got: %+v)", expected1, hs1)
	}
	hs2 := configTests[2].config.collectScopedPairHistory(&IdentifierGroup{scope: "sample", keys: []string{"foo", "bar"}})
	expected2 := []string{"Test1,\\ world!, test1,\\ test", "Test2,\\ world!, test2,\\ test"}
	if !reflect.DeepEqual(hs2, expected2) {
		t.Errorf("collectScopedPairHistory incorrect (expected: %+v, got: %+v)", expected2, hs2)
	}
	hs3 := configTests[2].config.collectScopedPairHistory(&IdentifierGroup{scope: "sample", keys: []string{"foo", "bar", "baz"}})
	expected3 := []string{"Test1,\\ world!, test1,\\ test, baz"}
	if !reflect.DeepEqual(hs3, expected3) {
		t.Errorf("collectScopedPairHistory incorrect (expected: %+v, got: %+v)", expected3, hs3)
	}
	hs4 := configTests[2].config.collectScopedPairHistory(&IdentifierGroup{scope: "foo", keys: []string{"test"}})
	expected4 := []string{}
	if !reflect.DeepEqual(hs4, expected4) {
		t.Errorf("collectScopedPairHistory incorrect (expected: %+v, got: %+v)", expected4, hs4)
	}
}
