package loading

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/steps/choice"
	"github.com/Samasource/jen/internal/specification/steps/do"
	"github.com/Samasource/jen/internal/specification/steps/execute"
	"github.com/Samasource/jen/internal/specification/steps/input"
	"github.com/Samasource/jen/internal/specification/steps/option"
	"github.com/Samasource/jen/internal/specification/steps/options"
	"github.com/Samasource/jen/internal/specification/steps/render"
	"github.com/kylelemons/go-gypsy/yaml"
	"sort"
)

func LoadSpec(node yaml.Map) (*specification.Spec, error) {
	spec := new(specification.Spec)

	// Load metadata
	metadata, err := getRequiredMap(node, "metadata")
	if err != nil {
		return nil, err
	}
	spec.Name, err = getRequiredString(metadata, "Name")
	if err != nil {
		return nil, err
	}
	spec.Description, err = getRequiredString(metadata, "description")
	if err != nil {
		return nil, err
	}
	spec.Version, err = getRequiredString(metadata, "version")
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

func loadActions(node yaml.Map) ([]specification.Action, error) {
	var actions []specification.Action
	for name, value := range node {
		stepList, ok := value.(yaml.List)
		if !ok {
			return nil, fmt.Errorf("value of action %q must be a list", name)
		}
		steps, err := loadSteps(stepList)
		if err != nil {
			return nil, fmt.Errorf("failed to load action %q: %w", name, err)
		}
		actions = append(actions, specification.Action{Name: name, Steps: steps})
	}
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Name < actions[j].Name
	})
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
			name: "render",
			fct:  loadRenderStep,
		},
		{
			name: "exec",
			fct:  loadExecStep,
		},
		{
			name: "do",
			fct:  loadDoStep,
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

func loadInputStep(node yaml.Map, ifCondition string) (executable.Executable, error) {
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

func loadOptionStep(node yaml.Map, ifCondition string) (executable.Executable, error) {
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

func loadOptionsStep(node yaml.Map, ifCondition string) (executable.Executable, error) {
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

func loadChoiceStep(node yaml.Map, ifCondition string) (executable.Executable, error) {
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

func loadRenderStep(node yaml.Map, ifCondition string) (executable.Executable, error) {
	source, err := getRequiredString(node, "source")
	if err != nil {
		return nil, err
	}

	return render.Render{
		If:     ifCondition,
		Source: source,
	}, nil
}

func loadExecStep(node yaml.Map, ifCondition string) (executable.Executable, error) {
	command, err := getRequiredString(node, "command")
	if err != nil {
		return nil, err
	}

	return execute.Execute{
		If:      ifCondition,
		Command: command,
	}, nil
}

func loadDoStep(node yaml.Map, ifCondition string) (executable.Executable, error) {
	action, err := getRequiredString(node, "action")
	if err != nil {
		return nil, err
	}

	return do.Do{
		If:     ifCondition,
		Action: action,
	}, nil
}
