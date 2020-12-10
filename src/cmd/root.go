package cmd

import (
	"fmt"
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
	c.PersistentPreRunE = func(*cobra.Command, []string) error {
		return initialize(config)
	}
	return c
}

func initialize(config *model.Config) error {
	jenHomeDir, ok := os.LookupEnv("JEN_HOME")
	if !ok {
		jenHomeDir = "~/.jen"
	}
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("finding jen home dir: %v", err)
	}
	Log("Using jen home dir: %s", jenHomeDir)
	config.JenDir = strings.ReplaceAll(jenHomeDir, "~", home)
	config.TemplatesDir = path.Join(config.JenDir, TemplatesDirName)
	return nil
}
