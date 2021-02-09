package persist

import (
	"fmt"
	"strings"

	"github.com/kylelemons/go-gypsy/yaml"
)

func getRequiredMap(node yaml.Map, key string) (yaml.Map, error) {
	child, ok := node[key]
	if !ok {
		return nil, fmt.Errorf("missing required property %q", key)
	}
	m, ok := child.(yaml.Map)
	if !ok {
		// WORKAROUND: go-gypsy lib incorrectly loads "{}" empty object as a literal string
		str, _ := getString(child)
		if str == "{}" {
			return yaml.Map{}, nil
		}
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

// getOptionalMapOrRawStringOrRawStrings retrieves the child map with given key or, if child is a raw string, it returns a map with
// raw string stored in a property keyed with defaultSubKey. This is to support steps that have two alternate syntaxes,
// a long-hand syntax using a map with multiple properties and a short-hand syntax with a raw string that specifies
// only the value of defaultSubKey. If defaultSubKey is an empty string, then only the long-hand map syntax is tried.
func getOptionalMapOrRawStringOrRawStrings(node yaml.Node, key, defaultSubKey string) (yaml.Map, bool, error) {
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
		str, ok := getString(child)
		if ok {
			_map = yaml.Map{
				defaultSubKey: yaml.Scalar(str),
			}
			return _map, true, nil
		}

		// Try list of raw strings
		list := getOptionalListOfScalar(child)
		if list != nil {
			_map = yaml.Map{
				defaultSubKey: list,
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

func getRequiredStringsOrStringFromMap(_map yaml.Map, key string) ([]string, error) {
	value, ok := _map[key]
	if !ok {
		return nil, fmt.Errorf("missing required property %q", key)
	}

	// If value is scalar, return a slice with just itself
	str, ok := getString(value)
	if ok {
		return []string{str}, nil
	}

	// Otherwise, value should be a list of raw strings
	list, ok := value.(yaml.List)
	if !ok {
		return nil, fmt.Errorf("property %q is expected to be either a raw string or a list of strings", key)
	}
	values := make([]string, 0, list.Len())
	for _, item := range list {
		str, ok := getString(item)
		if !ok {
			return nil, fmt.Errorf("property %q is expected to be either a raw string or a list of strings", key)
		}
		values = append(values, str)
	}
	return values, nil
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

func getStringInternal(_map yaml.Map, key string) (string, bool, error) {
	value, ok := _map[key]
	if !ok {
		return "", false, nil
	}
	str, ok := getString(value)
	if !ok {
		return "", false, fmt.Errorf("property %q must be a string", key)
	}
	return str, true, nil
}

func getString(node yaml.Node) (string, bool) {
	scalar, ok := node.(yaml.Scalar)
	if !ok {
		return "", false
	}
	str := scalar.String()
	// WORKAROUND: go-gypsy lib incorrectly loads `""` empty string as a literal of two double-quotes
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		return str[1 : len(str)-1], true
	}
	return str, true
}

func getOptionalListOfScalar(node yaml.Node) yaml.List {
	list, ok := node.(yaml.List)
	if !ok {
		return nil
	}

	// Ensure all list children are scalars
	for _, item := range list {
		if _, ok := item.(yaml.Scalar); !ok {
			return nil
		}
	}

	return list
}

func getOptionalBool(_map yaml.Map, key string, defaultValue bool) (bool, error) {
	value, ok, err := getStringInternal(_map, key)
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
