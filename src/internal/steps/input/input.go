package input

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
)

// Prompt represents a single text input user prompt
type Prompt struct {
	Message string
	Var     string
	Default string
}

// Execute prompts user for input value
func (p Prompt) Execute(config *model.Config) error {
	// Is var overriden?
	_, ok := config.VarOverrides[p.Var]
	if ok {
		return nil
	}

	// Compute message
	message, err := evaluation.EvalPromptValueTemplate(config.Values, config.BinDirs, p.Message)
	if err != nil {
		return err
	}

	// Compute default value
	defaultValue, ok := config.Values.Variables[p.Var]
	if !ok {
		defaultValue, err = evaluation.EvalPromptValueTemplate(config.Values, config.BinDirs, p.Default)
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

	config.Values.Variables[p.Var] = value
	return config.OnValuesChanged()
}
