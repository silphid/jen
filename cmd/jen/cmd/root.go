package cmd

import (
	"github.com/silphid/jen/cmd/jen/cmd/do"
	"github.com/silphid/jen/cmd/jen/cmd/exec"
	"github.com/silphid/jen/cmd/jen/cmd/export"
	"github.com/silphid/jen/cmd/jen/cmd/internal"
	"github.com/silphid/jen/cmd/jen/cmd/list"
	"github.com/silphid/jen/cmd/jen/cmd/pull"
	"github.com/silphid/jen/cmd/jen/cmd/require"
	"github.com/silphid/jen/cmd/jen/cmd/shell"
	"github.com/silphid/jen/cmd/jen/cmd/versioning"
	"github.com/silphid/jen/cmd/jen/internal/logging"
	"github.com/spf13/cobra"
)

// NewRoot creates the root cobra command
func NewRoot(version string) *cobra.Command {
	c := &cobra.Command{
		Use:   "jen",
		Short: "Jen is a code generator and script runner for creating and maintaining projects",
		Long: `Jen is a code generator and script runner that simplifies prompting for values, creating a new project
from those values and a given template, registering the project with your cloud infrastructure and CI/CD, and then
continues to support you throughout development in executing project-related commands and scripts using the same values.`,
		SilenceUsage: true,
	}

	var options internal.Options
	c.PersistentFlags().BoolVarP(&logging.Verbose, "verbose", "v", false, "display verbose messages")
	c.PersistentFlags().StringVarP(&options.TemplateName, "template", "t", "", "Name of template to use (defaults to prompting user)")
	c.PersistentFlags().BoolVarP(&options.SkipConfirm, "yes", "y", false, "skip all confirmation prompts")
	c.PersistentFlags().StringSliceVarP(&options.VarOverrides, "set", "s", []string{}, "sets a project variable manually (can be used multiple times)")
	c.AddCommand(versioning.New(version))
	c.AddCommand(pull.New())
	c.AddCommand(do.New(&options))
	c.AddCommand(exec.New(&options))
	c.AddCommand(shell.New(&options))
	c.AddCommand(list.New(&options))
	c.AddCommand(export.New(&options))
	c.AddCommand(require.New(&options))
	return c
}
