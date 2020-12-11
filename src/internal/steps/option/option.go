package option

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
	"strconv"
)

type Prompt struct {
	Message string
	Var     string
	Default bool
}

func (p Prompt) String() string {
	return "option"
}

func (p Prompt) Execute(config *model.Config) error {
	// Show prompt
	message, err := evaluation.EvalPromptValueTemplate(config.Values, p.Message)
	if err != nil {
		return err
	}
	prompt := &survey.Confirm{
		Message: message,
		Default: p.Default,
	}
	value := false
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	config.Values.Variables[p.Var] = strconv.FormatBool(value)
	return config.OnValuesChanged()
}
