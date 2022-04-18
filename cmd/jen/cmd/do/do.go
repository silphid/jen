package do

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/silphid/jen/cmd/jen/cmd/internal"
	"github.com/silphid/jen/cmd/jen/internal/exec"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "do",
		Short: "Executes an action from a template's spec.yaml",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			return run(options, args)
		},
	}
}

func run(options *internal.Options, optionalActionName []string) error {
	execContext, err := options.NewContext()
	if err != nil {
		return err
	}

	// If action name not specified, prompt user to select it from list of available actions
	actionName := ""
	if len(optionalActionName) == 0 {
		actionName, err = promptAction(execContext)
		if err != nil {
			return err
		}
	} else {
		actionName = optionalActionName[0]
	}

	// Retrieve action by name
	action := execContext.GetAction(actionName)
	if action == nil {
		return fmt.Errorf("action %q not found in spec file", actionName)
	}

	return action.Execute(execContext)
}

func promptAction(context exec.Context) (string, error) {
	actions := context.GetActionNames()
	prompt := &survey.Select{
		Message: "Select action to execute",
		Options: actions,
	}
	var actionName string
	if err := survey.AskOne(prompt, &actionName); err != nil {
		return "", err
	}
	return actionName, nil
}
