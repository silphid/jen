package specification

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/prompts/choice"
	"github.com/Samasource/jen/internal/specification/prompts/input"
	"github.com/Samasource/jen/internal/specification/prompts/option"
	"github.com/Samasource/jen/internal/specification/prompts/options"
	"github.com/kylelemons/go-gypsy/yaml"
	"strings"
)

func LoadActions(node yaml.Map) ([]Action, error) {
	var actions []Action
	for name, value := range node {
		stepList, ok := value.(yaml.List)
		if !ok {
			return nil, fmt.Errorf("value of action %q must be a list", name)
		}
		steps, err := loadSteps(stepList)
		if err != nil {
			return nil, fmt.Errorf("failed to load action %q: %w", name, err)
		}
		actions = append(actions, Action{Name: name, Steps: steps})
	}
	return actions, nil
}

func loadSteps(list yaml.List) ([]executable.Executable, error) {
	var steps []executable.Executable
	for idx, value := range list {
		stepMap, ok := value.(yaml.Map)
		if !ok {
			return nil, fmt.Errorf("value of step #%d must be an object", idx+1)
		}
		step, err := loadStep(stepMap)
		if err != nil {
			return nil, fmt.Errorf("failed to load step #%d: %w", idx+1, err)
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func loadStep(node yaml.Map) (executable.Executable, error) {
	ifCondition, err := getOptionalString(node, "if", "")
	if err != nil {
		return nil, err
	}

	items := []struct {
		name string
		fct  func(node yaml.Map, ifCondition string) (executable.Executable, error)
	}{
		{
			name: "input",
			fct:  loadTextPrompt,
		},
		{
			name: "option",
			fct:  loadOptionPrompt,
		},
		{
			name: "options",
			fct:  loadOptionsPrompt,
		},
		{
			name: "choice",
			fct:  loadChoicePrompt,
		},
	}

	for _, x := range items {
		child, ok, err := getOptionalMap(node, x.name)
		if err != nil {
			return nil, err
		}
		if ok {
			return x.fct(child, ifCondition)
		}
	}

	return nil, fmt.Errorf("unknown step type")
}

func loadTextPrompt(node yaml.Map, ifCondition string) (executable.Executable, error) {
	question, err := getRequiredString(node, "question")
	if err != nil {
		return nil, err
	}
	variable, err := getRequiredString(node, "var")
	if err != nil {
		return nil, err
	}
	defaultValue, err := getOptionalString(node, "default", "")
	if err != nil {
		return nil, err
	}
	return input.Prompt{
		If:       ifCondition,
		Question: question,
		Var:      variable,
		Default:  defaultValue,
	}, nil
}

func loadOptionPrompt(node yaml.Map, ifCondition string) (executable.Executable, error) {
	question, err := getRequiredString(node, "question")
	if err != nil {
		return nil, err
	}
	variable, err := getRequiredString(node, "var")
	if err != nil {
		return nil, err
	}
	defaultValue, err := getOptionalBool(node, "default", false)
	if err != nil {
		return nil, err
	}
	return option.Prompt{
		If:       ifCondition,
		Question: question,
		Var:      variable,
		Default:  defaultValue,
	}, nil
}

func loadOptionsPrompt(node yaml.Map, ifCondition string) (executable.Executable, error) {
	question, err := getRequiredString(node, "question")
	if err != nil {
		return nil, err
	}

	// Load children
	list, err := getRequiredList(node, "items")
	if err != nil {
		return nil, err
	}
	var items []options.Item
	for _, child := range list {
		childMap, ok := child.(yaml.Map)
		if !ok {
			return nil, fmt.Errorf("items of %q property must be objects", "options")
		}
		question, err := getRequiredString(childMap, "text")
		if err != nil {
			return nil, err
		}
		variable, err := getRequiredString(childMap, "var")
		if err != nil {
			return nil, err
		}
		defaultValue, err := getOptionalBool(childMap, "default", false)
		if err != nil {
			return nil, err
		}
		items = append(items, options.Item{
			Text:    question,
			Var:     variable,
			Default: defaultValue,
		})
	}

	return options.Prompt{
		If:       ifCondition,
		Question: question,
		Items:    items,
	}, nil
}

func loadChoicePrompt(node yaml.Map, ifCondition string) (executable.Executable, error) {
	question, err := getRequiredString(node, "question")
	if err != nil {
		return nil, err
	}
	variable, err := getRequiredString(node, "var")
	if err != nil {
		return nil, err
	}
	defaultValue, err := getOptionalString(node, "default", "")
	if err != nil {
		return nil, err
	}

	// Load children
	list, err := getRequiredList(node, "items")
	if err != nil {
		return nil, err
	}
	var items []choice.Item
	for _, child := range list {
		childMap, ok := child.(yaml.Map)
		if !ok {
			return nil, fmt.Errorf("items of %q property must be objects", "options")
		}
		text, err := getRequiredString(childMap, "text")
		if err != nil {
			return nil, err
		}
		value, err := getRequiredString(childMap, "value")
		if err != nil {
			return nil, err
		}
		items = append(items, choice.Item{
			Text:  text,
			Value: value,
		})
	}

	return choice.Prompt{
		If:       ifCondition,
		Question: question,
		Var:      variable,
		Default:  defaultValue,
		Items:    items,
	}, nil
}

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
