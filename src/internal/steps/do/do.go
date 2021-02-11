package do

import (
	"fmt"

	"github.com/Samasource/jen/src/internal/exec"
)

// Do represents a reference to another action within same spec file to which
// execution will be delegated
type Do struct {
	Action string
}

func (d Do) String() string {
	return "do"
}

// Execute executes another action with given name within same spec file
func (d Do) Execute(context exec.Context) error {
	action := context.GetAction(d.Action)
	if action == nil {
		return fmt.Errorf("action %q not found for do step", d.Action)
	}
	return action.Execute(context)
}
