package option

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/specification"
)

type Prompt struct {
	Message string
	Var     string
	Default bool
}

func (p Prompt) String() string {
	return "option"
}

func (p Prompt) Execute(context specification.Context) error {
	// Show prompt
	prompt := &survey.Confirm{
		Message: p.Message,
		Default: p.Default,
	}
	value := false
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	context.Values.Variables[p.Var] = value
	return nil
}
