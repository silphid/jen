package versioning

import (
	"fmt"

	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Displays build version",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
