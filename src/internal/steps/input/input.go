package input

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
)

type Prompt struct {
	Message string
	Var     string
	Default string
}

func (p Prompt) Execute(config *model.Config) error {
	// Show prompt
	message, err := evaluation.EvalPromptValueTemplate(config.Values, config.PathEnvVar, p.Message)
	if err != nil {
		return err
	}
	defaultValue, err := evaluation.EvalPromptValueTemplate(config.Values, config.PathEnvVar, p.Default)
	if err != nil {
		return err
	}
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
	}
	value := ""
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	config.Values.Variables[p.Var] = value
	return config.OnValuesChanged()
}
