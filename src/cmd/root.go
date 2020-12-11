package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/cmd/exec"
	"github.com/Samasource/jen/internal/persist"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Samasource/jen/cmd/do"
	. "github.com/Samasource/jen/internal/constant"
	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func NewRoot(config *model.Config) *cobra.Command {
	c := &cobra.Command{
		Use:   "jen",
		Short: "Jen is a code generator and script runner for creating and maintaining projects",
		Long: `Jen is a code generator and script runner that simplifies prompting for values, creating a new project
from those values and a given template, registering the project with your cloud infrastructure and CI/CD, and then
continues to support you throughout development in executing project-related commands and scripts using the same values.`,
		SilenceUsage: true,
	}

	c.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "display verbose messages")
	c.PersistentFlags().StringVarP(&config.TemplateName, "template", "t", "", "Name of template to use (defaults to prompting user)")
	c.PersistentFlags().BoolVarP(&config.SkipConfirm, "yes", "y", false, "skip all confirmation prompts")
	c.AddCommand(do.New(config))
	c.AddCommand(exec.New(config))
	c.PersistentPreRunE = func(*cobra.Command, []string) error {
		return initialize(config)
	}
	return c
}

func initialize(config *model.Config) error {
	var err error
	config.JenDir, err = getJenHomeDir()
	if err != nil {
		return err
	}
	config.TemplatesDir = path.Join(config.JenDir, TemplatesDirName)
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
		config.OnValuesChanged()
	}

	config.TemplateDir = path.Join(config.TemplatesDir, config.TemplateName)
	config.Spec, err = persist.LoadSpecFromDir(config.TemplateDir)
	return err
}

func getJenHomeDir() (string, error) {
	jenHomeDir, ok := os.LookupEnv("JEN_HOME")
	if !ok {
		jenHomeDir = "~/.jen"
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("finding jen home dir: %v", err)
	}
	Log("Using jen home dir: %s", jenHomeDir)
	return strings.ReplaceAll(jenHomeDir, "~", home), nil
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
