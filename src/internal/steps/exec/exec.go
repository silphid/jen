package exec

import (
	"path/filepath"

	"github.com/Samasource/jen/internal/model"
	"github.com/Samasource/jen/internal/shell"
)

type Exec struct {
	Commands []string
}

func (e Exec) String() string {
	return "exec"
}

func (e Exec) Execute(config *model.Config) error {
	dir, err := filepath.Abs(config.ProjectDir)
	if err != nil {
		return err
	}

	return shell.Execute(config.Values.Variables, dir, config.PathEnvVar, e.Commands...)
}
