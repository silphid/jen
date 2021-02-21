package choice

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/exec"
)

// Item represent one of the multiple choices prompted to user
type Item struct {
	Text  string
	Value string
}

// Prompt represents a user prompt for a single choice among many
type Prompt struct {
	Message string
	Var     string
	Default string
	Items   []Item
}

func (p Prompt) String() string {
	return "choice"
}

// Execute prompts user for choice value
func (p Prompt) Execute(context exec.Context) error {
	// Is var already set manually?
	if context.IsVarOverriden(p.Var) {
		return nil
	}

	// Collect option texts and find default index
	defaultIndex := 0
	currentValue, _ := context.GetVars()[p.Var]
	var options []string
	for i, item := range p.Items {
		text, err := evaluation.EvalTemplate(context, item.Text)
		if err != nil {
			return err
		}
		options = append(options, text)

		// Is this item the current value?
		if item.Value == currentValue {
			defaultIndex = i
		}
	}

	// Show prompt
	message, err := evaluation.EvalTemplate(context, p.Message)
	if err != nil {
		return err
	}
	prompt := &survey.Select{
		Message: message,
		Options: options,
		Default: defaultIndex,
	}
	var value int
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	vars := context.GetVars()
	vars[p.Var] = p.Items[value].Value
	return context.SetVars(vars)
}
