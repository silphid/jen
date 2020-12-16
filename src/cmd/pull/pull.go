package pull

import (
	"github.com/Samasource/jen/internal/model"
	"github.com/Samasource/jen/internal/shell"
	"github.com/spf13/cobra"
)

func New(config *model.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pulls latest templates from git repo",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return run(config)
		},
	}
}

func run(config *model.Config) error {
	return shell.Execute(nil, config.JenDir, "git pull")
}
