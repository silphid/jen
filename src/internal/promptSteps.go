/*package internal

import (
	"github.com/AlecAivazis/survey/v2"
)

func (step *StringStep) Execute(context Context) error {
	prompt := &survey.Input{
		Message: step.Title,
	}

	value := ""
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	context.Values[step.Name] = value
	return nil
}

func (step *SecretStep) Execute(context Context) error {
	prompt := &survey.Password{
		Message: step.Title,
	}

	value := ""
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	context.Values[step.Name] = value
	return nil
}

func (step *OptionStep) Execute(context Context) error {
	prompt := &survey.Confirm{
		Message: step.Title,
	}

	value := false
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	context.Values[step.Name] = value
	return nil
}

func (step *MultiOptionStep) Execute(context Context) error {
	var titles []string
	for _, item := range step.Items {
		titles = append(titles, item.Title)
	}

	prompt := &survey.MultiSelect{
		Message: step.Title,
		Options: titles,
	}

	var value []int
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	for _, index := range value {
		name := step.Items[index].Name
		context.Values[name] = true
	}
	return nil
}

func (step *SelectStep) Execute(context Context) error {
	var titles []string
	for _, item := range step.Items {
		titles = append(titles, item.Title)
	}

	prompt := &survey.Select{
		Message: step.Title,
		Options: titles,
	}

	var value int
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}

	context.Values[step.Name] = step.Items[value].Value
	return nil
}
*/
