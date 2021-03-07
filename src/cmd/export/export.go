package export

import (
	"fmt"
	"sort"

	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: `Outputs vars in "export VAR=value" format to be sourced using "$(jen export)"`,
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

	vars := execContext.GetShellVars(false)
	sort.Strings(vars)
	for _, v := range vars {
		fmt.Printf("export %s\n", v)
	}
	return nil
}
