package choice

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
)

type Item struct {
	Text  string
	Value string
}

type Prompt struct {
	Message string
	Var     string
	Default string
	Items   []Item
}

func (p Prompt) String() string {
	return "choice"
}

func (p Prompt) Execute(config *model.Config) error {
	// Collect option texts
	var options []string
	for _, item := range p.Items {
		text, err := evaluation.EvalPromptValueTemplate(config.Values, item.Text)
		if err != nil {
			return err
		}
		options = append(options, text)
	}

	// Show prompt
	message, err := evaluation.EvalPromptValueTemplate(config.Values, p.Message)
	if err != nil {
		return err
	}
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	var value int
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	config.Values.Variables[p.Var] = p.Items[value].Value
	return config.OnValuesChanged()
}
