package list

import (
	"strings"

	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/Samasource/jen/src/cmd/list/actions"
	"github.com/Samasource/jen/src/cmd/list/scripts"
	"github.com/Samasource/jen/src/cmd/list/templates"
	"github.com/Samasource/jen/src/cmd/list/vars"
	"github.com/Samasource/jen/src/internal/shell"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists available templates, actions, variables or scripts",
		RunE: func(_ *cobra.Command, args []string) error {
			return run(options, args)
		},
	}
	c.AddCommand(actions.New(options))
	c.AddCommand(scripts.New(options))
	c.AddCommand(templates.New(options))
	c.AddCommand(vars.New(options))
	return c
}

func run(options *internal.Options, args []string) error {
	execContext, err := options.NewContext()
	if err != nil {
		return err
	}

	return shell.Execute(execContext.GetShellVars(), "", strings.Join(args, " "))
}
