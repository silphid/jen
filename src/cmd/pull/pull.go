package pull

import (
	"github.com/silphid/jen/src/internal/home"
	"github.com/silphid/jen/src/internal/shell"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pulls latest templates from git repo",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			jenHome, err := home.GetOrCloneRepo()
			if err != nil {
				return err
			}

			return shell.Execute(nil, jenHome, "git pull")
		},
	}
}
