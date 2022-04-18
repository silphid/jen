package list

import (
	"github.com/silphid/jen/cmd/jen/cmd/internal"
	"github.com/silphid/jen/cmd/jen/cmd/list/actions"
	"github.com/silphid/jen/cmd/jen/cmd/list/scripts"
	"github.com/silphid/jen/cmd/jen/cmd/list/templates"
	"github.com/silphid/jen/cmd/jen/cmd/list/vars"
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
