package persist

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/src/internal/constant"
	"github.com/Samasource/jen/src/internal/home"
	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/project"
)

// LoadOrCreateJenFile loads the current project's jen file and, if it doesn't
// exists, it prompts users whether to create it.
func LoadOrCreateJenFile(config *model.Config) error {
	projectDir, err := project.GetProjectDir()
	if err != nil {
		return err
	}
	if projectDir == "" {
		if !config.SkipConfirm {
			err := confirmCreateJenFile()
			if err != nil {
				return err
			}
		}
		err := SaveConfig(config)
		if err != nil {
			return err
		}
	}

	err = LoadConfig(config)
	if err != nil {
		return err
	}

	jenHomeDir, err := home.GetJenHomeDir()
	if err != nil {
		return err
	}
	templatesDir := path.Join(jenHomeDir, constant.TemplatesDirName)

	if config.TemplateName == "" {
		config.TemplateName, err = promptTemplate(templatesDir)
		if err != nil {
			return fmt.Errorf("prompting for template: %w", err)
		}
		config.OnValuesChanged()
	}

	// Apply command-line variable overrides
	if len(config.VarOverrides) > 0 {
		for key, value := range config.VarOverrides {
			config.Values.Variables[key] = value
		}
		config.OnValuesChanged()
	}

	config.TemplateDir = path.Join(templatesDir, config.TemplateName)
	config.Spec, err = LoadSpecFromDir(config.TemplateDir)
	if err != nil {
		return err
	}

	config.BinDirs = []string{
		path.Join(jenHomeDir, "bin"),
		path.Join(config.TemplateDir, "bin"),
	}
	return nil
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
