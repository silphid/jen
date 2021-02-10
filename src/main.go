package main

import (
	"os"

	"github.com/Samasource/jen/src/cmd"
	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/persist"
	"github.com/Samasource/jen/src/internal/project"
)

func main() {
	config := &model.Config{}
	config.OnValuesChanged = func() error {
		projectDir, err := project.GetProjectDir()
		if err != nil {
			return err
		}
		return persist.SaveConfig(config, projectDir)
	}

	rootCmd := cmd.NewRoot(config)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
