package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
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
		return nil, fmt.Errorf("invalid JSON in config file: %v", err)
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

func (config *Config) collectHistory(id *Identifier) []string {
	values := []string{}
	if _, ok := config.Scopes[id.scope]; ok {
		for _, value := range config.Scopes[id.scope].Values {
			if v, ok := value[id.key]; ok {
				values = append(values, v)
			}
		}
	}
	return values
}

func (config *Config) collectScopedPairHistory(idg *IdentifierGroup) []string {
	values := []string{}
	if _, ok := config.Scopes[idg.scope]; ok {
		for _, value := range config.Scopes[idg.scope].Values {
			contained := true
			var vs []string
			for _, key := range idg.keys {
				if v, ok := value[key]; ok {
					vs = append(vs, strings.Replace(v, ", ", ",\\ ", -1))
				} else {
					contained = false
					break
				}
			}
			if contained {
				values = append(values, strings.Join(vs, ", "))
			}
		}
	}
	return values
}
