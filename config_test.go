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
							"test": "Test",
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
          "test": "Test"
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
