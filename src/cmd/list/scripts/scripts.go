package scripts

import (
	"fmt"

	"github.com/silphid/jen/src/cmd/internal"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:     "scripts",
		Aliases: []string{"script"},
		Short:   "Lists scripts available in current template (including shared scripts)",
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

	scripts, err := execContext.GetScripts()
	if err != nil {
		return err
	}
	for _, script := range scripts {
		fmt.Println(script)
	}
	return nil
}
