package require

import (
	"fmt"
	"os"

	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "require",
		Short: "Validates that all given variables are defined in current environment",
		RunE: func(_ *cobra.Command, args []string) error {
			return run(options, args)
		},
	}
}

func run(options *internal.Options, args []string) error {
	valid := true
	for _, arg := range args {
		_, ok := os.LookupEnv(arg)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing required variable %q\n", arg)
			valid = false
		}
	}

	if !valid {
		os.Exit(1)
	}
	return nil
}
