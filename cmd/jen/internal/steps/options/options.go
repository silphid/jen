package options

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/silphid/jen/cmd/jen/internal/evaluation"
	"github.com/silphid/jen/cmd/jen/internal/exec"
	"github.com/silphid/jen/cmd/jen/internal/helpers/variables"
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
	var defaultIndices []int
	var options []string
	for i, item := range p.Items {
		// Compute message
		text, err := evaluation.EvalTemplate(context, item.Text)
		if err != nil {
			return err
		}
		options = append(options, text)

		// Compute default value
		defaultValue, ok := variables.TryGetBool(vars, item.Var)
		if !ok {
			defaultValue = item.Default
		}
		if defaultValue {
			defaultIndices = append(defaultIndices, i)
		}
	}

	// Show prompt
	message, err := evaluation.EvalTemplate(context, p.Message)
	if err != nil {
		return err
	}
	var indices []int
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Default: defaultIndices,
	}
	if err := survey.AskOne(prompt, &indices); err != nil {
		return err
	}

	// Clear all options
	for _, item := range p.Items {
		vars[item.Var] = false
	}

	// Enable selected options
	for _, index := range indices {
		name := p.Items[index].Var
		vars[name] = true
	}

	return context.SetVars(vars)
}
