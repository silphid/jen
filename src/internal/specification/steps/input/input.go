package input

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/model"
)

type Prompt struct {
	Message string
	Var     string
	Default string
}

func (p Prompt) Execute(config model.Config) error {
	// Show prompt
	prompt := &survey.Input{
		Message: p.Message,
		Default: p.Default,
	}
	value := ""
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	config.Values.Variables[p.Var] = value
	return config.SaveJenFile()
}
