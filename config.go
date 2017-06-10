package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"sort"
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

func stringifyValue(values map[string]string) string {
	keys := make([]string, len(values))
	i := 0
	for k := range values {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	buf := new(bytes.Buffer)
	for _, k := range keys {
		buf.WriteByte('\u0001')
		buf.WriteString(k)
		buf.WriteByte('\u0002')
		buf.WriteString(values[k])
	}
	return buf.String()
}

func (config *Config) collectHistory(id *Identifier) []string {
	var values []string
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
	var values []string
	if _, ok := config.Scopes[idg.scope]; ok {
		for _, value := range config.Scopes[idg.scope].Values {
			contained := true
			var vs []string
			for _, key := range idg.keys {
				if v, ok := value[key]; ok {
					if strings.Contains(v, ", ") {
						v = strings.Replace(v, ", ", ",\\ ", -1)
					}
					vs = append(vs, v)
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
