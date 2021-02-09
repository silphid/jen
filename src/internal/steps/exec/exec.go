package exec

import (
	"path/filepath"

	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/shell"
)

// Exec represent a set of shell commands
type Exec struct {
	Commands []string
}

func (e Exec) String() string {
	return "exec"
}

// Execute executes one or multiple shell commands with project's variables and bin dirs
func (e Exec) Execute(config *model.Config) error {
	dir, err := filepath.Abs(config.ProjectDir)
	if err != nil {
		return err
	}

	return shell.Execute(config.Values.Variables, dir, config.BinDirs, e.Commands...)
}
