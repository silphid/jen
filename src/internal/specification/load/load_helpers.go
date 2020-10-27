package load

import (
	"fmt"
	"github.com/kylelemons/go-gypsy/yaml"
	"strings"
)

func getOptionalMap(node yaml.Map, key string) (yaml.Map, bool, error) {
	child, ok := node[key]
	if !ok {
		return nil, false, nil
	}
	m, ok := child.(yaml.Map)
	if !ok {
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

func getRequiredString(node yaml.Map, key string) (string, error) {
	value, ok, err := getString(node, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("missing required property %q", key)
	}
	return value, nil
}

func getOptionalString(node yaml.Map, key string, defaultValue string) (string, error) {
	value, ok, err := getString(node, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return defaultValue, nil
	}
	return value, nil
}

func getString(node yaml.Map, key string) (string, bool, error) {
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
	value, ok, err := getString(node, key)
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
