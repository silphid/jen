package do

import (
	"fmt"

	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/persist"
	"github.com/spf13/cobra"
)

// New creates the "jen do" cobra sub-command
func New(config *model.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "do",
		Short: "Executes an action from a template's spec.yaml",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return run(config, args[0])
		},
	}
}

func run(config *model.Config, actionName string) error {
	err := persist.LoadOrCreateJenFile(config)
	if err != nil {
		return err
	}

	action, ok := config.Spec.Actions[actionName]
	if !ok {
		return fmt.Errorf("action %q not found in spec file", actionName)
	}

	return action.Execute(config)
}
