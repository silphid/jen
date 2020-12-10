package main

import (
	"github.com/Samasource/jen/cmd"
	"github.com/Samasource/jen/internal/model"
	"github.com/Samasource/jen/internal/persist"
	"os"
)

func main() {
	config := &model.Config{}
	config.OnValuesChanged = func() error {
		return persist.SaveJenFile(config)
	}

	rootCmd := cmd.NewRoot(config)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
