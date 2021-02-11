package options

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/exec"
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
func (p Prompt) Execute(context exec.Context) error {
	// Are all vars overriden?
	allVarsOverriden := true
	for _, item := range p.Items {
		if !context.IsVarOverriden(item.Var) {
			allVarsOverriden = false
			break
		}
	}
	if allVarsOverriden {
		return nil
	}

	vars := context.GetVars()

	// Collect option texts and default values
	var indices []int
	var options []string
	for i, item := range p.Items {
		// Compute message
		text, err := evaluation.EvalPromptValueTemplate(context, item.Text)
		if err != nil {
			return err
		}
		options = append(options, text)

		// Compute default value
		defaultValue := item.Default
		defaultString, ok := vars[item.Var]
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
	message, err := evaluation.EvalPromptValueTemplate(context, p.Message)
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
	for _, item := range p.Items {
		vars[item.Var] = "false"
	}

	// Enable selected options
	for _, index := range indices {
		name := p.Items[index].Var
		vars[name] = "true"
	}
	return context.SaveProject()
}
