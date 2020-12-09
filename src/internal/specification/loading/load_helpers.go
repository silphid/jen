package loading

import (
	"fmt"
	"github.com/kylelemons/go-gypsy/yaml"
	"strings"
)

func getRequiredMap(node yaml.Map, key string) (yaml.Map, error) {
	child, ok := node[key]
	if !ok {
		return nil, fmt.Errorf("missing required property %q", key)
	}
	m, ok := child.(yaml.Map)
	if !ok {
		return nil, fmt.Errorf("property %q must be an object", key)
	}
	return m, nil
}

func getOptionalMap(node yaml.Node, key string) (yaml.Map, bool, error) {
	_map, ok := node.(yaml.Map)
	if !ok {
		return nil, false, nil
	}
	child, ok := _map[key]
	if !ok {
		return nil, false, nil
	}
	m, ok := child.(yaml.Map)
	if !ok {
		return nil, false, fmt.Errorf("property %q must be an object", key)
	}
	return m, true, nil
}

// getOptionalMapOrRawString retrieves the child map with given key or, if child is a raw string, it returns a map with
// raw string stored in a property keyed with defaultSubKey. This is to support steps that have two alternate syntaxes,
// a long-hand syntax using a map with multiple properties and a short-hand syntax with a raw string that specifies
// only the value of defaultSubKey. If defaultSubKey is an empty string, then only the long-hand map syntax is tried.
func getOptionalMapOrRawString(node yaml.Node, key, defaultSubKey string) (yaml.Map, bool, error) {
	_map, ok := node.(yaml.Map)
	if !ok {
		return nil, false, nil
	}
	child, ok := _map[key]
	if !ok {
		return nil, false, nil
	}

	if defaultSubKey != "" {
		// Try raw string
		scalar, ok := child.(yaml.Scalar)
		if ok {
			_map = yaml.Map{
				defaultSubKey: scalar,
			}
			return _map, true, nil
		}
	}

	// Try map
	m, ok := child.(yaml.Map)
	if !ok {
		if defaultSubKey != "" {
			return nil, false, fmt.Errorf("property %q must be an object or raw string", key)
		}
		return nil, false, fmt.Errorf("property %q must be an object", key)
	}

	return m, true, nil
}

func getRequiredList(node yaml.Map, key string) (yaml.List, error) {
	child, ok := node[key]
	if !ok {
		return nil, fmt.Errorf("missing required property %q", key)
	}
	list, ok := child.(yaml.List)
	if !ok {
		return nil, fmt.Errorf("property %q must be a list", key)
	}
	return list, nil
}

func getRequiredStringFromMap(node yaml.Node, key string) (string, error) {
	_map, ok := node.(yaml.Map)
	if !ok {
		return "", fmt.Errorf("expected object")
	}
	value, ok, err := getStringInternal(_map, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("missing required property %q", key)
	}
	return value, nil
}

func getOptionalStringFromMap(node yaml.Node, key string, defaultValue string) (string, error) {
	_map, ok := node.(yaml.Map)
	if !ok {
		return defaultValue, nil
	}
	value, ok, err := getStringInternal(_map, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return defaultValue, nil
	}
	return value, nil
}

func getStringInternal(node yaml.Map, key string) (string, bool, error) {
	value, ok := node[key]
	if !ok {
		return "", false, nil
	}
	scalar, ok := value.(yaml.Scalar)
	if !ok {
		return "", false, fmt.Errorf("property %q must be a string", key)
	}
	return scalar.String(), true, nil
}

func getOptionalBool(node yaml.Map, key string, defaultValue bool) (bool, error) {
	value, ok, err := getStringInternal(node, key)
	if err != nil {
		return false, err
	}
	if !ok {
		return defaultValue, nil
	}
	switch strings.ToLower(value) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool value: %q", value)
	}
}
