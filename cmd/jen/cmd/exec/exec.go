package exec

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/silphid/jen/src/cmd/internal"
	"github.com/silphid/jen/src/internal/exec"
	"github.com/silphid/jen/src/internal/shell"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Executes custom scripts or arbitrary shell commands with project's environment variables",
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

	// If no args specified, prompt user to select from list of available custom scripts
	if len(args) == 0 {
		script, err := promptScript(execContext)
		if err != nil {
			return err
		}
		args = []string{script}
	}

	return shell.Execute(execContext.GetShellVars(true), "", strings.Join(args, " "))
}

func promptScript(context exec.Context) (string, error) {
	scripts, err := context.GetScripts()
	if err != nil {
		return "", err
	}
	prompt := &survey.Select{
		Message: "Select script to execute",
		Options: scripts,
	}
	var script string
	if err := survey.AskOne(prompt, &script); err != nil {
		return "", err
	}
	return script, nil
}
