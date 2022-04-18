package variables

import "github.com/silphid/jen/cmd/jen/internal/helpers/conversion"

// TryGetString tries to fetch given variable from map and return it as a string,
// also returning whether the variable was successfully found and converted.
func TryGetString(vars map[string]interface{}, name string) (string, bool) {
	value, ok := vars[name]
	if !ok {
		return "", false
	}

	str, err := conversion.ToString(value)
	if err != nil {
		return "", false
	}
	return str, true
}

// TryGetBool tries to fetch given variable from map and return it as a bool,
// also returning whether the variable was successfully found and converted.
func TryGetBool(vars map[string]interface{}, name string) (bool, bool) {
	value, ok := vars[name]
	if !ok {
		return false, false
	}

	b, err := conversion.ToBool(value)
	if err != nil {
		return false, false
	}
	return b, true
}
