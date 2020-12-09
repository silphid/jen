package main

import (
	"github.com/Samasource/jen/cmd"
	"github.com/Samasource/jen/internal/model"
	"os"
)

func main() {
	config := &model.Config{}
	rootCmd := cmd.NewRoot(config)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
