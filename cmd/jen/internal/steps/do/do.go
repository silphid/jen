package do

import (
	"fmt"

	"github.com/silphid/jen/src/internal/exec"
)

// Do represents a reference to another action within same spec file to which
// execution will be delegated
type Do struct {
	Actions []string
}

func (d Do) String() string {
	return "do"
}

// Execute executes another action with given name within same spec file
func (d Do) Execute(context exec.Context) error {
	for _, action := range d.Actions {
		action := context.GetAction(action)
		if action == nil {
			return fmt.Errorf("action %q not found for do step", action)
		}
		err := action.Execute(context)
		if err != nil {
			return err
		}
	}
	return nil
}
