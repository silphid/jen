package exec

import (
	"strings"

	"github.com/Samasource/jen/internal/persist"
	"github.com/Samasource/jen/internal/shell"

	"github.com/Samasource/jen/internal/model"
	"github.com/spf13/cobra"
)

// New creates the "jen exec" cobra sub-command
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
	err := persist.LoadOrCreateJenFile(config)
	if err != nil {
		return err
	}

	return shell.Execute(config.Values.Variables, "", config.BinDirs, strings.Join(args, " "))
}
