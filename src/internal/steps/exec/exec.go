package exec

import (
	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/project"
	"github.com/Samasource/jen/src/internal/shell"
)

// Exec represent a set of shell commands
type Exec struct {
	Commands []string
}

func (e Exec) String() string {
	return "exec"
}

// Execute runs one or multiple shell commands with project's variables and bin dirs
func (e Exec) Execute(config *model.Config) error {
	projectDir, err := project.GetDir()
	if err != nil {
		return err
	}

	return shell.Execute(config.Values.Variables, projectDir, config.BinDirs, e.Commands...)
}
