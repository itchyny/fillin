package main

import "strings"

// Config for fillin
type Config struct {
	Scopes map[string]*Scope `json:"scopes"`
}

// Scope holds pairs of values
type Scope struct {
	Values []map[string]string `json:"values"`
}

func (config *Config) collectHistory(id *Identifier) []string {
	values := []string{}
	added := make(map[string]bool)
	if _, ok := config.Scopes[id.scope]; ok {
		for _, value := range config.Scopes[id.scope].Values {
			if v, ok := value[id.key]; ok && !added[v] {
				values = append(values, v)
				added[v] = true
			}
		}
	}
	return values
}

func (config *Config) collectScopedPairHistory(idg *IdentifierGroup) []string {
	values := []string{}
	added := make(map[string]bool)
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
				v := strings.Join(vs, ", ")
				if !added[v] {
					values = append(values, v)
					added[v] = true
				}
			}
		}
	}
	return values
}
