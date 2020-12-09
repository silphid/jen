package choice

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/specification"
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

func (p Prompt) Execute(context specification.Context) error {
	// Collect option texts
	var options []string
	for _, item := range p.Items {
		options = append(options, item.Text)
	}

	// Show prompt
	prompt := &survey.Select{
		Message: p.Message,
		Options: options,
	}
	var value int
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	context.Values.Variables[p.Var] = p.Items[value].Value
	return nil
}
