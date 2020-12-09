package cmd

import (
	"fmt"
	"github.com/Samasource/jen/cmd/do"
	. "github.com/Samasource/jen/internal/constant"
	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
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
	c.AddCommand(do.New(config))

	c.PersistentPreRunE = func(*cobra.Command, []string) error {
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
	return c
}
