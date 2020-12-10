package options

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/model"
)

type Item struct {
	Text    string
	Var     string
	Default bool
}

type Prompt struct {
	Message string
	Items   []Item
}

func (p Prompt) String() string {
	return "options"
}

func (p Prompt) Execute(config model.Config) error {
	// Collect option texts
	var options []string
	for _, item := range p.Items {
		options = append(options, item.Text)
	}

	// Show prompt
	prompt := &survey.MultiSelect{
		Message: p.Message,
		Options: options,
	}
	var indices []int
	if err := survey.AskOne(prompt, &indices); err != nil {
		return err
	}

	// Clear all options
	for i := range p.Items {
		name := p.Items[i].Var
		config.Values.Variables[name] = "false"
	}

	// Enable selected options
	for _, index := range indices {
		name := p.Items[index].Var
		config.Values.Variables[name] = "true"
	}
	return config.SaveJenFile()
}
