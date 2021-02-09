package choice

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
)

// Item represent one of the multiple choices prompted to user
type Item struct {
	Text  string
	Value string
}

// Prompt represents a user prompt for a single choice among many
type Prompt struct {
	Message string
	Var     string
	Default string
	Items   []Item
}

func (p Prompt) String() string {
	return "choice"
}

// Execute prompts user for choice value
func (p Prompt) Execute(config *model.Config) error {
	// Is var already set manually?
	_, ok := config.SetVars[p.Var]
	if ok {
		return nil
	}

	// Collect option texts
	var options []string
	for _, item := range p.Items {
		text, err := evaluation.EvalPromptValueTemplate(config.Values, config.PathEnvVar, item.Text)
		if err != nil {
			return err
		}
		options = append(options, text)
	}

	// Show prompt
	message, err := evaluation.EvalPromptValueTemplate(config.Values, config.PathEnvVar, p.Message)
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
