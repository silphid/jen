package specification

import (
	"fmt"
	"github.com/Samasource/jen/internal"
)

type Action struct {
	Name  string
	Steps []Executable
}

func (a Action) String() string {
	return a.Name
}

func (a Action) Execute(context Context) error {
	internal.Log("Executing action %q", a.Name)
	for _, step := range a.Steps {
		if err := step.Execute(context); err != nil {
			return fmt.Errorf("failed to execute action %q, step %q: %w", a.Name, step.(fmt.Stringer).String(), err)
		}
	}
	return nil
}
