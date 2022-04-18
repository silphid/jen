package exec

import (
	"github.com/silphid/jen/cmd/jen/internal/exec"
	"github.com/silphid/jen/cmd/jen/internal/shell"
)

// Exec represent a set of shell commands
type Exec struct {
	Commands []string
}

func (e Exec) String() string {
	return "exec"
}

// Execute runs one or multiple shell commands with project's variables and bin dirs
func (e Exec) Execute(context exec.Context) error {
	return shell.Execute(context.GetShellVars(true), context.GetProjectDir(), e.Commands...)
}
