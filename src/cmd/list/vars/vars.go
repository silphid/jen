package vars

import (
	"fmt"
	"sort"

	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:     "vars",
		Aliases: []string{"var"},
		Short:   "Lists variables defined in current project and their current values",
		Args:    cobra.NoArgs,
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

	vars := execContext.GetVars()

	// Sort names
	names := make([]string, 0, len(vars))
	for name := range vars {
		names = append(names, name)
	}
	sort.Strings(names)

	// Print names and values
	for _, name := range names {
		value := vars[name]
		fmt.Printf("%s: %v\n", name, value)
	}
	return nil
}
