package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Samasource/jen/src/cmd/do"
	"github.com/Samasource/jen/src/cmd/exec"
	"github.com/Samasource/jen/src/cmd/pull"
	"github.com/Samasource/jen/src/internal/constant"
	"github.com/Samasource/jen/src/internal/helpers"
	"github.com/Samasource/jen/src/internal/logging"
	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/shell"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type flags struct {
	templateName string
	skipConfirm  bool
	varOverrides []string
}

// NewRoot creates the root cobra command
func NewRoot(config *model.Config) *cobra.Command {
	c := &cobra.Command{
		Use:   "jen",
		Short: "Jen is a code generator and script runner for creating and maintaining projects",
		Long: `Jen is a code generator and script runner that simplifies prompting for values, creating a new project
from those values and a given template, registering the project with your cloud infrastructure and CI/CD, and then
continues to support you throughout development in executing project-related commands and scripts using the same values.`,
		SilenceUsage: true,
	}

	var flags flags
	c.PersistentFlags().BoolVarP(&logging.Verbose, "verbose", "v", false, "display verbose messages")
	c.PersistentFlags().StringVarP(&flags.templateName, "template", "t", "", "Name of template to use (defaults to prompting user)")
	c.PersistentFlags().BoolVarP(&flags.skipConfirm, "yes", "y", false, "skip all confirmation prompts")
	c.PersistentFlags().StringSliceVarP(&flags.varOverrides, "set", "s", []string{}, "sets a project variable manually (can be used multiple times)")
	c.AddCommand(pull.New(config))
	c.AddCommand(do.New(config))
	c.AddCommand(exec.New(config))
	c.PersistentPreRunE = func(*cobra.Command, []string) error {
		return initialize(config, flags)
	}
	return c
}

func initialize(config *model.Config, flags flags) error {
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

	config.TemplatesDir = path.Join(config.JenDir, constant.TemplatesDirName)
	config.ProjectDir, err = findProjectDirUpFromWorkDir()
	if err != nil {
		return err
	}
	config.VarOverrides, err = parseOverrideVars(flags.varOverrides)
	config.TemplateName = flags.templateName
	config.SkipConfirm = flags.skipConfirm
	return err
}

func cloneJenRepo(jenHomeDir, jenRepo string) error {
	// Jen dir already exists and is a valid git working copy?
	homeExists := helpers.PathExists(jenHomeDir)
	if homeExists {
		dotGitDir := path.Join(jenHomeDir, ".git")
		if helpers.PathExists(dotGitDir) {
			// Jen dir is a valid git repo
			return nil
		}

		// Not a valid git repo, therefore must be empty, so we can clone into it
		infos, err := ioutil.ReadDir(jenHomeDir)
		if err != nil {
			return fmt.Errorf("listing content of jen dir %q to ensure it's empty before cloning into it: %w", jenHomeDir, err)
		}
		if len(infos) > 0 {
			return fmt.Errorf("jen dir %q already exists, is not a valid git working copy and already contains files so we cannot clone into it (please delete or empty it)", jenHomeDir)
		}
	}

	// Clone jen repo
	logging.Log("Cloning jen templates repo %q into jen dir %q", jenRepo, jenHomeDir)
	return shell.Execute(nil, "", nil, fmt.Sprintf("git clone %s %s", jenRepo, jenHomeDir))
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
	logging.Log("Using jen home dir: %s", jenHomeDir)
	return strings.ReplaceAll(jenHomeDir, "~", home), nil
}

func findProjectDirUpFromWorkDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding project's root dir: %w", err)
	}

	for {
		filePath := path.Join(dir, constant.JenFileName)
		if helpers.PathExists(filePath) {
			return dir, nil
		}
		if dir == "/" {
			return "", nil
		}
		dir = path.Dir(dir)
	}
}

var varOverrideRegexp = regexp.MustCompile(`^(\w+)=(.*)$`)

func parseOverrideVars(rawVarOverrides []string) (map[string]string, error) {
	varOverrides := make(map[string]string, len(rawVarOverrides))
	for _, raw := range rawVarOverrides {
		submatch := varOverrideRegexp.FindStringSubmatch(raw)
		if submatch == nil {
			return nil, fmt.Errorf("failed to parse set variable %q", raw)
		}
		varOverrides[submatch[1]] = submatch[2]
	}
	return varOverrides, nil
}
