package option

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
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
func (p Prompt) Execute(config *model.Config) error {
	// Is var already set manually?
	_, ok := config.SetVars[p.Var]
	if ok {
		return nil
	}

	// Compute default value
	defaultValue := p.Default
	defaultString, ok := config.Values.Variables[p.Var]
	if ok {
		var err error
		defaultValue, err = strconv.ParseBool(defaultString)
		if err != nil {
			return err
		}
	}

	// Show prompt
	message, err := evaluation.EvalPromptValueTemplate(config.Values, config.PathEnvVar, p.Message)
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

	config.Values.Variables[p.Var] = strconv.FormatBool(value)
	return config.OnValuesChanged()
}
