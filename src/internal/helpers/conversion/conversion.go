package conversion

import (
	"fmt"
	"strconv"
)

// ToString converts an abstract value (can be either string or bool) to its string representation
func ToString(value interface{}) (string, error) {
	strValue, ok := value.(string)
	if ok {
		return strValue, nil
	}

	boolValue, ok := value.(bool)
	if ok {
		return strconv.FormatBool(boolValue), nil
	}

	return "", fmt.Errorf("failed to convert type into string: %t", value)
}

// ToBool converts an abstract value (can be either string or bool) to its string representation
func ToBool(value interface{}) (bool, error) {
	boolValue, ok := value.(bool)
	if ok {
		return boolValue, nil
	}

	strValue, ok := value.(string)
	if ok {
		return strconv.ParseBool(strValue)
	}

	return false, fmt.Errorf("failed to convert type into bool: %t", value)
}
