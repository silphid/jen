package pull

import (
	"github.com/Samasource/jen/src/internal/home"
	"github.com/Samasource/jen/src/internal/shell"
	"github.com/spf13/cobra"
)

// New creates the "jen pull" cobra sub-command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pulls latest templates from git repo",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			jenHome, err := home.CloneJenRepo()
			if err != nil {
				return err
			}

			return shell.Execute(nil, jenHome, nil, "git pull")
		},
	}
}
