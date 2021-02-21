package option

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/exec"
	"github.com/Samasource/jen/src/internal/helpers/variables"
)

// Prompt represents a boolean user prompt
type Prompt struct {
	Message string
	Var     string
	Default bool
}

func (p Prompt) String() string {
	return "option"
}

// Execute prompts user for a boolean value
func (p Prompt) Execute(context exec.Context) error {
	if context.IsVarOverriden(p.Var) {
		return nil
	}

	vars := context.GetVars()

	// Compute default value
	defaultValue, ok := variables.TryGetBool(vars, p.Var)
	if !ok {
		defaultValue = p.Default
	}

	// Show prompt
	message, err := evaluation.EvalTemplate(context, p.Message)
	if err != nil {
		return err
	}
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
	}
	value := false
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	vars[p.Var] = value
	return context.SetVars(vars)
}
