package do

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	. "github.com/Samasource/jen/internal/constant"
	"github.com/Samasource/jen/internal/model"
	"github.com/Samasource/jen/internal/specification/loading"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
)

func New(config *model.Config) *cobra.Command {
	c := &cobra.Command{
		Use:   "do",
		Short: "Executes an action from a template's spec.yaml",
		Args:  cobra.ExactArgs(1),
		RunE: func(*cobra.Command, []string) error {
			// TODO: Do action
			return nil
		},
	}

	c.PersistentFlags().StringVarP(&config.TemplateName, "template", "t", "", "Name of template to use (defaults to prompting user)")

	c.PersistentPreRunE = func(*cobra.Command, []string) error {
		var err error
		config.ProjectDir, err = findProjectDir()
		if err != nil {
			return err
		}
		if config.ProjectDir == "" {
			// TODO: Project not initialized
		} else {
			// TODO: Load jen.yaml values
		}

		if config.TemplateName == "" {
			config.TemplateName, err = promptTemplate(config.TemplatesDir)
			if err != nil {
				return err
			}
			config.TemplateDir = path.Join(config.TemplatesDir)
			config.Spec, err = loading.LoadSpecFromDir(config.TemplateDir)
		}
		return nil
	}
	return c
}

func findProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding project's root dir: %w", err)
	}

	for {
		filePath := path.Join(dir, ValuesFileName)
		exists, err := fileExists(filePath)
		if err != nil {
			return "", fmt.Errorf("finding project's root dir: %w", err)
		}
		if exists {
			return dir, nil
		}
		dir = path.Dir(dir)
		if dir == "" {
			return "", nil
		}
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
		return "", err
	}

	// Build list of choices
	var templates []string
	var titles []string
	for _, info := range infos {
		template := info.Name()
		templateDir := path.Join(templatesDir, template)
		spec, err := loading.LoadSpecFromDir(templateDir)
		if err == nil {
			templates = append(templates, template)
			titles = append(titles, fmt.Sprintf("%s - %s", template, spec.Description))
		}
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
