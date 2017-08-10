package main

import (
	"fmt"
	"os"
	"strings"
)

// Identifier ...
type Identifier struct {
	scope string
	key   string
	defaultValue string
}

func (id *Identifier) prompt() string {
	if id.scope == "" {
		if len(id.defaultValue) > 0 {
			return fmt.Sprintf("%s(%s): ", id.key, id.defaultValue)
		}
		return fmt.Sprintf("%s: ", id.key)
	}
	if len(id.defaultValue) > 0 {
		return fmt.Sprintf("[%s] %s(%s): ", id.scope, id.key, id.defaultValue)
	}
	return fmt.Sprintf("[%s] %s: ", id.scope, id.key)
}

// IdentifierGroup ...
type IdentifierGroup struct {
	scope string
	keys  []string
	defaultValues []string
}

func (idg *IdentifierGroup) prompt() string {
	var messages []string
	for index, key := range idg.keys {
		var msg = key
		if len(idg.defaultValues[index]) > 0 {
			msg = msg + fmt.Sprintf("(%s)", idg.defaultValues[index])
		}
		messages = append(messages, msg)
	}
	return fmt.Sprintf("[%s] %s: ", idg.scope, strings.Join(messages, ", "))
}

func defaultValue(key string, scope string) string {
	value, ok := os.LookupEnv(defaultValueKey(key, scope))
	if ok && len(value) > 0 {
		return value
	}
	return ""
}

func defaultValueKey(key string, scope string) string {
	if len(scope) > 0 {
		return fmt.Sprintf("%s_%s", scope, key)
	} else {
		return key
	}
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
	var keys, defaultValues []string
	added := make(map[string]bool)
	for _, id := range identifiers {
		if scope == id.scope && !added[id.key] {
			keys = append(keys, id.key)
			defaultValues = append(defaultValues, defaultValue(id.key, scope))
			added[id.key] = true
		}
	}
	return &IdentifierGroup{scope: scope, keys: keys, defaultValues: defaultValues}
}

func insert(values map[string]map[string]string, id *Identifier, value string) {
	if _, ok := values[id.scope]; !ok {
		values[id.scope] = make(map[string]string)
	}
	values[id.scope][id.key] = value
}

func empty(values map[string]map[string]string) bool {
	for scope := range values {
		for key := range values[scope] {
			if values[scope][key] != "" {
				return false
			}
		}
	}
	return true
}

func lookup(values map[string]map[string]string, id *Identifier) string {
	if v, ok := values[id.scope]; ok {
		if v, ok := v[id.key]; ok {
			return v
		}
	}
	return ""
}
