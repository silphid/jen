package persist

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal/model"
	"io/ioutil"
	"path"
	"strings"
)

func LoadOrCreateJenFile(config *model.Config) error {
	if config.ProjectDir == "" {
		if !config.SkipConfirm {
			err := confirmCreateJenFile()
			if err != nil {
				return err
			}
		}
		err := SaveJenFile(config)
		if err != nil {
			return err
		}
	}

	err := LoadJenFile(config)
	if err != nil {
		return err
	}

	if config.TemplateName == "" {
		config.TemplateName, err = promptTemplate(config.TemplatesDir)
		if err != nil {
			return fmt.Errorf("prompting for template: %w", err)
		}
		config.OnValuesChanged()
	}

	config.TemplateDir = path.Join(config.TemplatesDir, config.TemplateName)
	config.Spec, err = LoadSpecFromDir(config.TemplateDir)
	return err
}

func confirmCreateJenFile() error {
	var result bool
	err := survey.AskOne(&survey.Confirm{
		Message: "Jen project not found. Do you want to initialize current directory as your project root?",
		Default: false,
	}, &result)
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("cancelled by user")
	}
	return nil
}

func promptTemplate(templatesDir string) (string, error) {
	// Read templates dir
	infos, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		return "", fmt.Errorf("reading templates directory %q: %w", templatesDir, err)
	}

	// Build list of choices
	var templates []string
	var titles []string
	for _, info := range infos {
		template := info.Name()
		if strings.HasPrefix(template, ".") {
			continue
		}
		templateDir := path.Join(templatesDir, template)
		spec, err := LoadSpecFromDir(templateDir)
		if err != nil {
			return "", err
		}
		templates = append(templates, template)
		titles = append(titles, fmt.Sprintf("%s - %s", template, spec.Description))
	}

	// Any templates found?
	if len(templates) == 0 {
		return "", fmt.Errorf("no templates found in %q", templatesDir)
	}

	// Prompt
	prompt := &survey.Select{
		Message: "Select template",
		Options: titles,
	}
	var index int
	if err := survey.AskOne(prompt, &index); err != nil {
		return "", err
	}

	return templates[index], nil
}
