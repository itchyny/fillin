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
				"": {},
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
				"": {
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
				"sample": {
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
						{
							"foo": "Test1, world!",
							"bar": "test1, test",
							"qux": "qux",
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
        },
        {
          "bar": "test1, test",
          "foo": "Test1, world!",
          "qux": "qux"
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
	testCases := []struct {
		identifier *Identifier
		expected   []string
	}{
		{
			identifier: &Identifier{key: "foo"},
			expected:   []string{"Hello 1", "Hello 2"},
		},
		{
			identifier: &Identifier{key: "baz"},
			expected:   []string{"world!"},
		},
		{
			&Identifier{scope: "sample", key: "foo"},
			[]string{"Test1, world!", "Test2, world!"},
		},
		{
			&Identifier{scope: "foo", key: "test"},
			[]string{},
		},
	}
	for _, tc := range testCases {
		got := configTests[2].config.collectHistory(tc.identifier)
		if !reflect.DeepEqual(tc.expected, got) {
			t.Errorf("collectHistory incorrect (expected: %+v, got: %+v)", tc.expected, got)
		}
	}
}

func Test_collectScopedPairHistory(t *testing.T) {
	testCases := []struct {
		identifier *IdentifierGroup
		expected   []string
	}{
		{
			identifier: &IdentifierGroup{keys: []string{"foo", "baz"}},
			expected:   []string{"Hello 1, world!", "Hello 2, world!"},
		},
		{
			identifier: &IdentifierGroup{scope: "sample", keys: []string{"foo", "bar"}},
			expected:   []string{"Test1,\\ world!, test1,\\ test", "Test2,\\ world!, test2,\\ test"},
		},
		{
			identifier: &IdentifierGroup{scope: "sample", keys: []string{"foo", "bar", "baz"}},
			expected:   []string{"Test1,\\ world!, test1,\\ test, baz"},
		},
		{
			identifier: &IdentifierGroup{scope: "foo", keys: []string{"test"}},
			expected:   []string{},
		},
	}
	for _, tc := range testCases {
		got := configTests[2].config.collectScopedPairHistory(tc.identifier)
		if !reflect.DeepEqual(tc.expected, got) {
			t.Errorf("collectScopedPairHistory incorrect (expected: %+v, got: %+v)", tc.expected, got)
		}
	}
}
