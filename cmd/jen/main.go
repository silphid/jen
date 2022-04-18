package main

import (
	"os"

	"github.com/silphid/jen/cmd/jen/cmd"
)

var version string

func main() {
	rootCmd := cmd.NewRoot(version)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
