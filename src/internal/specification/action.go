package specification

import (
	"github.com/Samasource/jen/internal"
)

type Action struct {
	Name  string
	Steps Executables
}

type ActionMap map[string]Action

func (a Action) String() string {
	return a.Name
}

func (a Action) Execute(context Context) error {
	internal.Log("Executing sub-steps of action %q", a.Name)
	return a.Steps.Execute(context)
}
