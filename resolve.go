package main

import (
	"strings"
)

// Resolve asks the user to resolve the identifiers
func Resolve(identifiers []*Identifier, config *Config, p prompt) (map[string]map[string]string, error) {
	values := make(map[string]map[string]string)
	p.start()
	defer p.close()

	scopeAsked := make(map[string]bool)
	for _, id := range identifiers {
		if found(values, id) || id.scope == "" || scopeAsked[id.scope] {
			continue
		}
		scopeAsked[id.scope] = true
		idg := collect(identifiers, id.scope)
		if len(idg.keys) == 0 {
			continue
		}
		history := config.collectScopedPairHistory(idg)
		if len(history) == 0 {
			continue
		}
		p.setHistory(history)
		text, err := p.prompt(idg.prompt())
		if err != nil {
			return nil, err
		}
		xs := strings.Split(strings.TrimSuffix(text, "\n"), ", ")
		if len(xs) == len(idg.keys) {
			for i, key := range idg.keys {
				insert(values, &Identifier{scope: id.scope, key: key}, strings.Replace(xs[i], ",\\ ", ", ", -1))
			}
		}
	}

	for _, id := range identifiers {
		if found(values, id) {
			continue
		}
		p.setHistory(config.collectHistory(id))
		text, err := p.prompt(id.prompt())
		if err != nil {
			return nil, err
		}
		insert(values, id, strings.TrimSuffix(text, "\n"))
	}

	return values, nil
}
