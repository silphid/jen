package spec

import (
	"github.com/silphid/jen/cmd/jen/internal/exec"
	logging "github.com/silphid/jen/cmd/jen/internal/logging"
)

// Action represents a named executable that can be invoked from the
// command line via "jen do XXX" or via a "do" step
type Action struct {
	Name  string
	Steps exec.Executables
}

// ActionMap represents a dictionary mapping action names to their
// corresponding action
type ActionMap map[string]Action

func (a Action) String() string {
	return a.Name
}

// Execute executes many steps in sequence
func (a Action) Execute(context exec.Context) error {
	logging.Log("Executing action %q", a.Name)
	return a.Steps.Execute(context)
}
