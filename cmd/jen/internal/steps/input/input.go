package input

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/silphid/jen/cmd/jen/internal/evaluation"
	"github.com/silphid/jen/cmd/jen/internal/exec"
	"github.com/silphid/jen/cmd/jen/internal/helpers/variables"
)

// Prompt represents a single text input user prompt
type Prompt struct {
	Message string
	Var     string
	Default string
}

func (p Prompt) String() string {
	return "input"
}

// Execute prompts user for input value
func (p Prompt) Execute(context exec.Context) error {
	if context.IsVarOverriden(p.Var) {
		return nil
	}

	// Compute message
	message, err := evaluation.EvalTemplate(context, p.Message)
	if err != nil {
		return err
	}

	vars := context.GetVars()

	// Compute default value
	defaultValue, ok := variables.TryGetString(vars, p.Var)
	if !ok {
		defaultValue, err = evaluation.EvalTemplate(context, p.Default)
		if err != nil {
			return err
		}
	}

	// Show prompt
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
	}
	value := ""
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	vars[p.Var] = value
	return context.SetVars(vars)
}
