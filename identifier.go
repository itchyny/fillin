package main

// Identifier ...
type Identifier struct {
	scope string
	key   string
}

// IdentifierGroup ...
type IdentifierGroup struct {
	scope string
	keys  []string
}

func found(values map[string]map[string]string, id *Identifier) bool {
	if v, ok := values[id.scope]; ok {
		if _, ok := v[id.key]; ok {
			return true
		}
	}
	return false
}

func collect(identifiers []*Identifier, scope string) *IdentifierGroup {
	var keys []string
	added := make(map[string]bool)
	for _, id := range identifiers {
		if scope == id.scope && !added[id.key] {
			keys = append(keys, id.key)
			added[id.key] = true
		}
	}
	return &IdentifierGroup{scope: scope, keys: keys}
}

func insert(values map[string]map[string]string, id *Identifier, value string) {
	if _, ok := values[id.scope]; !ok {
		values[id.scope] = make(map[string]string)
	}
	values[id.scope][id.key] = value
}

func lookup(values map[string]map[string]string, id *Identifier) string {
	if v, ok := values[id.scope]; ok {
		if v, ok := v[id.key]; ok {
			return v
		}
	}
	return ""
}
