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
	"strings"
)

// NewRoot creates the root cobra command
func NewRoot(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "jen",
		Short: "Jen is a code generator and script runner for creating and maintaining projects",
		Long: `Jen is a code generator and script runner that simplifies prompting for values, creating a new project
from those values and a given template, registering the project with your cloud infrastructure and CI/CD, and then
continues to support you throughout development in executing project-related commands and scripts using the same values.`,
		SilenceUsage: true,
	}

	var options internal.Options

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Parse flags manually, because exec sub-command cannot parse them
		flags := getFlags(args)
		rootCmd.PersistentFlags().Parse(flags)
		return nil
	}

	rootCmd.PersistentFlags().BoolVarP(&logging.Verbose, "verbose", "v", false, "display verbose messages")
	rootCmd.PersistentFlags().StringVarP(&options.TemplateName, "template", "t", "", "Name of template to use (defaults to prompting user)")
	rootCmd.PersistentFlags().BoolVarP(&options.SkipConfirm, "yes", "y", false, "skip all confirmation prompts")
	rootCmd.PersistentFlags().BoolVar(&options.SkipPull, "skip-pull", false, "skip pulling latest template git repo")
	rootCmd.PersistentFlags().StringSliceVarP(&options.VarOverrides, "set", "s", []string{}, "sets a project variable manually (can be used multiple times)")

	rootCmd.AddCommand(versioning.New(version))
	rootCmd.AddCommand(pull.New())
	rootCmd.AddCommand(do.New(&options))
	rootCmd.AddCommand(exec.New(&options))
	rootCmd.AddCommand(list.New(&options))
	rootCmd.AddCommand(shell.New(&options))
	rootCmd.AddCommand(export.New(&options))
	rootCmd.AddCommand(require.New(&options))

	return rootCmd
}

// Split args into flags and sub-command
func getFlags(args []string) []string {
	var flagArgs []string
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			break
		}
		flagArgs = append(flagArgs, arg)
	}
	return flagArgs
}
