package cmd

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/Samasource/jen/src/cmd/do"
	"github.com/Samasource/jen/src/cmd/exec"
	"github.com/Samasource/jen/src/cmd/pull"
	"github.com/Samasource/jen/src/internal/constant"
	"github.com/Samasource/jen/src/internal/helpers"
	"github.com/Samasource/jen/src/internal/logging"
	"github.com/Samasource/jen/src/internal/model"
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
	c.AddCommand(pull.New())
	c.AddCommand(do.New(config))
	c.AddCommand(exec.New(config))
	c.PersistentPreRunE = func(*cobra.Command, []string) error {
		return initialize(config, flags)
	}
	return c
}

func initialize(config *model.Config, flags flags) error {
	var err error
	config.ProjectDir, err = findProjectDirUpFromWorkDir()
	if err != nil {
		return err
	}
	config.VarOverrides, err = parseOverrideVars(flags.varOverrides)
	config.TemplateName = flags.templateName
	config.SkipConfirm = flags.skipConfirm
	return err
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
