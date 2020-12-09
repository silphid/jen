package loading

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
	"github.com/Samasource/jen/internal/specification/steps/choice"
	"github.com/Samasource/jen/internal/specification/steps/do"
	"github.com/Samasource/jen/internal/specification/steps/exec"
	"github.com/Samasource/jen/internal/specification/steps/input"
	"github.com/Samasource/jen/internal/specification/steps/option"
	"github.com/Samasource/jen/internal/specification/steps/options"
	"github.com/Samasource/jen/internal/specification/steps/render"
	"github.com/kylelemons/go-gypsy/yaml"
)

func LoadSpec(node yaml.Map) (*specification.Spec, error) {
	spec := new(specification.Spec)

	// Load metadata
	metadata, err := getRequiredMap(node, "metadata")
	if err != nil {
		return nil, err
	}
	spec.Name, err = getRequiredStringFromMap(metadata, "Name")
	if err != nil {
		return nil, err
	}
	spec.Description, err = getRequiredStringFromMap(metadata, "description")
	if err != nil {
		return nil, err
	}
	spec.Version, err = getRequiredStringFromMap(metadata, "version")
	if err != nil {
		return nil, err
	}

	// Load actions
	actions, err := getRequiredMap(node, "actions")
	if err != nil {
		return nil, err
	}
	spec.Actions, err = loadActions(actions)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

func loadActions(node yaml.Map) (specification.ActionMap, error) {
	var actions []specification.Action
	for name, value := range node {
		stepList, ok := value.(yaml.List)
		if !ok {
			return nil, fmt.Errorf("value of action %q must be a list", name)
		}
		executables, err := loadExecutables(stepList)
		if err != nil {
			return nil, fmt.Errorf("failed to load action %q: %w", name, err)
		}
		actions = append(actions, specification.Action{Name: name, Steps: executables})
	}

	// Convert to map
	m := make(specification.ActionMap)
	for _, action := range actions {
		m[action.Name] = action
	}
	return m, nil
}

func loadExecutables(list yaml.List) (specification.Executables, error) {
	var executables specification.Executables
	for idx, value := range list {
		step, err := loadExecutable(value)
		if err != nil {
			return nil, fmt.Errorf("failed to load step #%d: %w", idx+1, err)
		}
		executables = append(executables, step)
	}
	return executables, nil
}

func loadExecutable(node yaml.Node) (specification.Executable, error) {
	// Special case for if step
	_map, ok := node.(yaml.Map)
	if ok {
		_, ok = _map["if"]
		if ok {
			return loadIfStep(_map)
		}
	}

	// Other steps
	items := []struct {
		name          string
		defaultSubKey string
		fct           func(node yaml.Map) (specification.Executable, error)
	}{
		{
			name: "input",
			fct:  loadInputStep,
		},
		{
			name: "option",
			fct:  loadOptionStep,
		},
		{
			name: "options",
			fct:  loadOptionsStep,
		},
		{
			name: "choice",
			fct:  loadChoiceStep,
		},
		{
			name:          "render",
			defaultSubKey: "source",
			fct:           loadRenderStep,
		},
		{
			name:          "exec",
			defaultSubKey: "commands",
			fct:           loadExecStep,
		},
		{
			name:          "do",
			defaultSubKey: "action",
			fct:           loadDoStep,
		},
	}

	for _, x := range items {
		_map, ok, err := getOptionalMapOrRawString(node, x.name, x.defaultSubKey)
		if err != nil {
			return nil, err
		}
		if ok {
			return x.fct(_map)
		}
	}

	return nil, fmt.Errorf("unknown step type")
}

func loadIfStep(_map yaml.Map) (specification.Executable, error) {
	condition, err := getRequiredStringFromMap(_map, "if")
	if err != nil {
		return nil, err
	}
	list, err := getRequiredList(_map, "then")
	if err != nil {
		return nil, err
	}
	executables, err := loadExecutables(list)
	if err != nil {
		return nil, err
	}
	return steps.If{
		Condition: condition,
		Then:      executables,
	}, nil
}

func loadInputStep(_map yaml.Map) (specification.Executable, error) {
	question, err := getRequiredStringFromMap(_map, "question")
	if err != nil {
		return nil, err
	}
	variable, err := getRequiredStringFromMap(_map, "var")
	if err != nil {
		return nil, err
	}
	defaultValue, err := getOptionalStringFromMap(_map, "default", "")
	if err != nil {
		return nil, err
	}
	return input.Prompt{
		Message: question,
		Var:     variable,
		Default: defaultValue,
	}, nil
}

func loadOptionStep(_map yaml.Map) (specification.Executable, error) {
	question, err := getRequiredStringFromMap(_map, "question")
	if err != nil {
		return nil, err
	}
	variable, err := getRequiredStringFromMap(_map, "var")
	if err != nil {
		return nil, err
	}
	defaultValue, err := getOptionalBool(_map, "default", false)
	if err != nil {
		return nil, err
	}
	return option.Prompt{
		Message: question,
		Var:     variable,
		Default: defaultValue,
	}, nil
}

func loadOptionsStep(_map yaml.Map) (specification.Executable, error) {
	question, err := getRequiredStringFromMap(_map, "question")
	if err != nil {
		return nil, err
	}

	// Load children
	list, err := getRequiredList(_map, "items")
	if err != nil {
		return nil, err
	}
	var items []options.Item
	for _, child := range list {
		childMap, ok := child.(yaml.Map)
		if !ok {
			return nil, fmt.Errorf("items of %q property must be objects", "options")
		}
		question, err := getRequiredStringFromMap(childMap, "text")
		if err != nil {
			return nil, err
		}
		variable, err := getRequiredStringFromMap(childMap, "var")
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
		Message: question,
		Items:   items,
	}, nil
}

func loadChoiceStep(_map yaml.Map) (specification.Executable, error) {
	question, err := getRequiredStringFromMap(_map, "question")
	if err != nil {
		return nil, err
	}
	variable, err := getRequiredStringFromMap(_map, "var")
	if err != nil {
		return nil, err
	}
	defaultValue, err := getOptionalStringFromMap(_map, "default", "")
	if err != nil {
		return nil, err
	}

	// Load children
	list, err := getRequiredList(_map, "items")
	if err != nil {
		return nil, err
	}
	var items []choice.Item
	for _, child := range list {
		childMap, ok := child.(yaml.Map)
		if !ok {
			return nil, fmt.Errorf("items of %q property must be objects", "options")
		}
		text, err := getRequiredStringFromMap(childMap, "text")
		if err != nil {
			return nil, err
		}
		value, err := getRequiredStringFromMap(childMap, "value")
		if err != nil {
			return nil, err
		}
		items = append(items, choice.Item{
			Text:  text,
			Value: value,
		})
	}

	return choice.Prompt{
		Message: question,
		Var:     variable,
		Default: defaultValue,
		Items:   items,
	}, nil
}

func loadRenderStep(_map yaml.Map) (specification.Executable, error) {
	source, err := getRequiredStringFromMap(_map, "source")
	if err != nil {
		return nil, err
	}

	return render.Render{
		Source: source,
	}, nil
}

func loadExecStep(_map yaml.Map) (specification.Executable, error) {
	commands, err := getRequiredStringsOrStringFromMap(_map, "commands")
	if err != nil {
		return nil, err
	}

	return exec.Exec{
		Commands: commands,
	}, nil
}

func loadDoStep(_map yaml.Map) (specification.Executable, error) {
	action, err := getRequiredStringFromMap(_map, "action")
	if err != nil {
		return nil, err
	}

	return do.Do{
		Action: action,
	}, nil
}
