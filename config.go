package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Config for fillin
type Config struct {
	Scopes map[string]*Scope `json:"scopes"`
}

// Scope holds pairs of values
type Scope struct {
	Values []map[string]string `json:"values"`
}

// ReadConfig loads a Config from reader
func ReadConfig(r io.Reader) (*Config, error) {
	var config Config
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return &config, nil
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// WriteConfig saves a Config to a writer
func WriteConfig(w io.Writer, config *Config) error {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}
