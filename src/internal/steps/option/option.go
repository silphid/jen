package option

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/exec"
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
	defaultValue := p.Default
	defaultString, ok := vars[p.Var]
	if ok {
		var err error
		defaultValue, err = strconv.ParseBool(defaultString)
		if err != nil {
			return err
		}
	}

	// Show prompt
	message, err := evaluation.EvalPromptValueTemplate(context, p.Message)
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

	vars[p.Var] = strconv.FormatBool(value)
	return context.SaveProject()
}
