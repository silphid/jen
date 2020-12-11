package exec

import (
	"github.com/Samasource/jen/internal/shell"
	"strings"

	"github.com/Samasource/jen/internal/model"
	"github.com/spf13/cobra"
)

func New(config *model.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Executes an arbitrary shell command with project's environment variables",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return run(config, args)
		},
	}
}

func run(config *model.Config, args []string) error {
	return shell.Execute(config.Values.Variables, "", strings.Join(args, " "))
}
