package options

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
)

// Item represent one of the multiple boolean values prompted to user
type Item struct {
	Text    string
	Var     string
	Default bool
}

// Prompt represents a user prompt for a set of individual boolean values
type Prompt struct {
	Message string
	Items   []Item
}

func (p Prompt) String() string {
	return "options"
}

// Execute prompts user for multiple individual boolean values
func (p Prompt) Execute(config *model.Config) error {
	// Are all vars overriden?
	allVarsOverriden := true
	for _, item := range p.Items {
		value, ok := config.VarOverrides[item.Var]
		if !ok {
			allVarsOverriden = false
			break
		}
		_, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("variable %q value %q failed to parse as boolean: %w", item.Var, value, err)
		}
	}
	if allVarsOverriden {
		return nil
	}

	// Collect option texts and default values
	var indices []int
	var options []string
	for i, item := range p.Items {
		// Compute message
		text, err := evaluation.EvalPromptValueTemplate(config.Values, config.BinDirs, item.Text)
		if err != nil {
			return err
		}
		options = append(options, text)

		// Compute default value
		defaultValue := item.Default
		defaultString, ok := config.Values.Variables[item.Var]
		if ok {
			var err error
			defaultValue, err = strconv.ParseBool(defaultString)
			if err != nil {
				return err
			}
		}
		if defaultValue {
			indices = append(indices, i)
		}
	}

	// Show prompt
	message, err := evaluation.EvalPromptValueTemplate(config.Values, config.BinDirs, p.Message)
	if err != nil {
		return err
	}
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Default: indices,
	}
	if err := survey.AskOne(prompt, &indices); err != nil {
		return err
	}

	// Clear all options
	for i := range p.Items {
		name := p.Items[i].Var
		config.Values.Variables[name] = "false"
	}

	// Enable selected options
	for _, index := range indices {
		name := p.Items[index].Var
		config.Values.Variables[name] = "true"
	}
	return config.OnValuesChanged()
}
