package specification

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/prompts"
	"github.com/Samasource/jen/internal/specification/prompts/option"
	"github.com/Samasource/jen/internal/specification/prompts/text"
	"github.com/kylelemons/go-gypsy/yaml"
	"strings"
)

func loadActions(node yaml.Map) ([]Action, error) {
	var actions []Action
	for name, value := range node {
		stepList, ok := value.(yaml.List)
		if !ok {
			return nil, fmt.Errorf("value of action %q must be a list", name)
		}
		steps, err := loadSteps(name, stepList)
		if err != nil {
			return nil, fmt.Errorf("failed to load action %q: %w", name, err)
		}
		actions = append(actions, Action{Name: name, Steps: steps})
	}
	return actions, nil
}

func loadSteps(list yaml.List) ([]*Executable, error) {
	var steps []*Executable
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

func loadStep(node yaml.Map) (*Executable, error) {
	ifCondition, err := getOptionalString(node, "if", "")
	if err != nil {
		return nil, err
	}
	if promptNode, ok := node["prompt"]; ok {
		promptMap, ok := promptNode.(yaml.Map)
		if !ok {
			return nil, fmt.Errorf("value of prompt must be an object")
		}
		return loadPrompt(promptMap, ifCondition)
	}
	return nil, fmt.Errorf("unknown step type")
}

func loadPrompt(node yaml.Map, ifCondition string) (Executable, error) {
	promptType, err := getOptionalString(node, "type", "text")
	if err != nil {
		return nil, err
	}
	switch promptType {
	case "text":
		return loadTextPrompt(node, ifCondition)
	case "option":
		return loadOptionPrompt(node, ifCondition)
	case "options":
		return loadOptionsPrompt(node, ifCondition)
	case "choice":
		return loadChoicePrompt(node, ifCondition)
	default:
		return nil, fmt.Errorf("unsupported prompt type %q", promptType)
	}
}

func loadTextPrompt(node yaml.Map, ifCondition string) (Executable, error) {
	x := text.Prompt{}
	x.If = ifCondition
	var err error
	x.Question, err = getRequiredString(node, "question")
	if err != nil {
		return nil, err
	}
	x.Var, err = getRequiredString(node, "var")
	if err != nil {
		return nil, err
	}
	x.Default, err = getOptionalString(node, "default", "")
	if err != nil {
		return nil, err
	}
	return x, nil
}

func loadOptionPrompt(node yaml.Map, ifCondition string) (Executable, error) {
	x := option.Prompt{}
	x.If = ifCondition
	var err error
	x.Question, err = getRequiredString(node, "question")
	if err != nil {
		return nil, err
	}
	x.Var, err = getRequiredString(node, "var")
	if err != nil {
		return nil, err
	}
	x.Default, err = getOptionalBool(node, "default", false)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func loadOptionsPrompt(node yaml.Map, ifCondition string) (Executable, error) {
	return nil, fmt.Errorf("not implemented")
}

func loadChoicePrompt(node yaml.Map, ifCondition string) (Executable, error) {
	return nil, fmt.Errorf("not implemented")
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
