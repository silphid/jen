package main

import (
	"os"

	"github.com/Samasource/jen/src/cmd"
	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/persist"
)

func main() {
	config := &model.Config{}
	config.OnValuesChanged = func() error {
		return persist.SaveConfig(config)
	}

	rootCmd := cmd.NewRoot(config)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
