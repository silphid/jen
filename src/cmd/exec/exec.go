package exec

import (
	"strings"

	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/Samasource/jen/src/internal/shell"
	"github.com/spf13/cobra"
)

// New creates the "jen exec" cobra sub-command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Executes an arbitrary shell command with project's environment variables",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return run(options, args)
		},
	}
}

func run(options *internal.Options, args []string) error {
	execContext, err := options.NewContext()
	if err != nil {
		return err
	}

	return shell.Execute(execContext.GetShellVars(), "", strings.Join(args, " "))
}
