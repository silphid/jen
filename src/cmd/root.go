package cmd

import (
	"fmt"
	"github.com/Samasource/jen/cmd/pull"
	"github.com/Samasource/jen/internal/shell"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Samasource/jen/cmd/do"
	"github.com/Samasource/jen/cmd/exec"
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
	c.AddCommand(pull.New(config))
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

	jenRepo, err := getJenRepo()
	if err != nil {
		return err
	}

	err = cloneJenRepo(config.JenDir, jenRepo)
	if err != nil {
		return err
	}

	config.TemplatesDir = path.Join(config.JenDir, TemplatesDirName)
	config.ProjectDir, err = findProjectDirUpFromWorkDir()
	return err
}

func cloneJenRepo(jenHomeDir, jenRepo string) error {
	// Jen dir already exists and is a valid git working copy?
	homeExists, err := pathExists(jenHomeDir)
	if err != nil {
		return err
	}
	if homeExists {
		dotGitDir := path.Join(jenHomeDir, ".git")
		dotGitExists, err := pathExists(dotGitDir)
		if err != nil {
			return err
		}
		if dotGitExists {
			// Jen dir is a valid git repo
			return nil
		} else {
			// Not a valid git repo, therefore must be empty, so we can clone into it
			infos, err := ioutil.ReadDir(jenHomeDir)
			if err != nil {
				return fmt.Errorf("listing content of jen dir %q to ensure it's empty before cloning into it: %w", jenHomeDir, err)
			}
			if len(infos) > 0 {
				return fmt.Errorf("jen dir %q already exists, is not a valid git working copy and already contains files so we cannot clone into it (please delete or empty it)", jenHomeDir)
			}
		}
	}

	// Clone jen repo
	Log("Cloning jen templates repo %q into jen dir %q", jenRepo, jenHomeDir)
	return shell.Execute(nil, "", fmt.Sprintf("git clone %s %s", jenRepo, jenHomeDir))
}

func getJenRepo() (string, error) {
	jenRepo, ok := os.LookupEnv("JEN_REPO")
	if !ok {
		return "", fmt.Errorf("please specify a JEN_REPO env var pointing to your jen templates git repo")
	}
	return jenRepo, nil
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

func findProjectDirUpFromWorkDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding project's root dir: %w", err)
	}

	for {
		filePath := path.Join(dir, JenFileName)
		exists, err := pathExists(filePath)
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

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("checking if %q path exists: %w", path, err)
	}
	return true, nil
}
