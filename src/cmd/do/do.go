package do

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	. "github.com/Samasource/jen/internal/constant"
	"github.com/Samasource/jen/internal/model"
	"github.com/Samasource/jen/internal/persist"
	"github.com/spf13/cobra"
)

func New(config *model.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "do",
		Short: "Executes an action from a template's spec.yaml",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return run(config, args[0])
		},
	}
}

func run(config *model.Config, actionName string) error {
	err := initialize(config)
	if err != nil {
		return fmt.Errorf("initializing config: %w", err)
	}

	action, ok := config.Spec.Actions[actionName]
	if !ok {
		return fmt.Errorf("action %q not found in spec file", actionName)
	}
	return action.Execute(config)
}

func initialize(config *model.Config) error {
	var err error
	config.ProjectDir, err = findProjectDirUpFromWorkDir()
	if err != nil {
		return err
	}

	err = loadOrCreateJenFile(config)
	if err != nil {
		return err
	}

	if config.TemplateName == "" {
		config.TemplateName, err = promptTemplate(config.TemplatesDir)
		if err != nil {
			return fmt.Errorf("prompting for template: %w", err)
		}
	}

	config.TemplateDir = path.Join(config.TemplatesDir, config.TemplateName)
	config.Spec, err = persist.LoadSpecFromDir(config.TemplateDir)
	return err
}

func loadOrCreateJenFile(config *model.Config) error {
	if config.ProjectDir == "" {
		if !config.SkipConfirm {
			err := confirmCreateJenFile()
			if err != nil {
				return err
			}
		}
		err := persist.SaveJenFile(config)
		if err != nil {
			return err
		}
	} else {
		err := persist.LoadJenFile(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func confirmCreateJenFile() error {
	var result bool
	err := survey.AskOne(&survey.Confirm{
		Message: "Do you want jen to initialize current directory as your project root?",
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

func findProjectDirUpFromWorkDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding project's root dir: %w", err)
	}

	for {
		filePath := path.Join(dir, JenFileName)
		exists, err := fileExists(filePath)
		if err != nil {
			return "", fmt.Errorf("finding project's root dir: %w", err)
		}
		if exists {
			return dir, nil
		}
		if dir == "/" {
			return "", nil
		}
		dir = path.Dir(dir)
	}
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("checking if %q file exists: %w", path, err)
	}
	return true, nil
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
		spec, err := persist.LoadSpecFromDir(templateDir)
		if err != nil {
			return "", nil
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
