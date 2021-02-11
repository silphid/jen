package do

import (
	"fmt"

	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/spf13/cobra"
)

// New creates the "jen do" cobra sub-command
func New(options internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "do",
		Short: "Executes an action from a template's spec.yaml",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return run(options, args[0])
		},
	}
}

func run(options internal.Options, actionName string) error {
	execContext, err := options.NewContext()
	if err != nil {
		return err
	}

	action := execContext.GetAction(actionName)
	if action == nil {
		return fmt.Errorf("action %q not found in spec file", actionName)
	}

	return action.Execute(execContext)
}
