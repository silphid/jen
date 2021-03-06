package shell

import (
	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/Samasource/jen/src/internal/shell"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Starts a sub-shell with project's environment variables and scripts in PATH",
		Args:  cobra.NoArgs,
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

	return shell.Execute(execContext.GetShellVars(), "", "$SHELL")
}
