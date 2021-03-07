package list

import (
	"github.com/Samasource/jen/src/cmd/internal"
	"github.com/Samasource/jen/src/cmd/list/actions"
	"github.com/Samasource/jen/src/cmd/list/scripts"
	"github.com/Samasource/jen/src/cmd/list/templates"
	"github.com/Samasource/jen/src/cmd/list/vars"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists available templates, actions, variables or scripts",
	}
	c.AddCommand(actions.New(options))
	c.AddCommand(scripts.New(options))
	c.AddCommand(templates.New(options))
	c.AddCommand(vars.New(options))
	return c
}
