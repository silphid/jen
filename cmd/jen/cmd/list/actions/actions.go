package actions

import (
	"fmt"

	"github.com/silphid/jen/cmd/jen/cmd/internal"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:     "actions",
		Aliases: []string{"action"},
		Short:   "Lists actions available in current template",
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

	for _, action := range execContext.GetActionNames() {
		fmt.Println(action)
	}
	return nil
}
